package controllers

import (
	"baggsy/backend/models"
	"baggsy/backend/utils"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterBagRequest represents the expected request body for registering a bag.
type RegisterBagRequest struct {
	QRCode      string `json:"qrCode"`                // QR code of the bag being registered
	BagType     string `json:"bagType"`               // "Parent" or "Child"
	ChildCount  int    `json:"childCount,omitempty"`  // For Parent bags, number of child bags expected
	ParentBagID *uint  `json:"parentBagID,omitempty"` // For Child bags, the ID of the parent bag
}

// RegisterBag handles the registration of a new parent or child bag.
func RegisterBag(c *gin.Context) {
	var req RegisterBagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If request body is not valid JSON or missing fields
		utils.HandleError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Validate BagType value
	if req.BagType != "Parent" && req.BagType != "Child" {
		utils.HandleError(c, http.StatusBadRequest, "BagType must be either 'Parent' or 'Child'")
		return
	}

	// Create a new Bag instance based on the request
	bag := models.Bag{
		QRCode:     req.QRCode,
		BagType:    req.BagType,
		ChildCount: 0,     // default, may override if parent
		Linked:     false, // default for new bag
		// LinkedToBill remains false by default
	}

	if req.BagType == "Parent" {
		// For parent bag, set the child count from request and ensure no ParentBagID is provided
		bag.ChildCount = req.ChildCount
		if req.ParentBagID != nil {
			utils.HandleError(c, http.StatusBadRequest, "Parent bag should not have a parentBagID")
			return
		}
	} else if req.BagType == "Child" {
		// For child bag, ParentBagID must be provided
		if req.ParentBagID == nil {
			utils.HandleError(c, http.StatusBadRequest, "Child bag must have a parentBagID")
			return
		}
		// Find the parent bag to link to
		var parent models.Bag
		if err := models.DB.First(&parent, *req.ParentBagID).Error; err != nil {
			utils.HandleError(c, http.StatusBadRequest, "Specified parent bag not found")
			return
		}
		if parent.BagType != "Parent" {
			utils.HandleError(c, http.StatusBadRequest, "Specified parent bag ID does not belong to a Parent bag")
			return
		}
		if parent.LinkedToBill {
			// Parent already linked to a bill, no new children can be added
			utils.HandleError(c, http.StatusBadRequest, "Cannot add child bags to a parent bag that is already linked to a bill")
			return
		}
		// Count existing children for this parent
		var currentChildrenCount int64
		models.DB.Model(&models.Bag{}).Where("parent_bag_id = ?", parent.ID).Count(&currentChildrenCount)
		if currentChildrenCount >= int64(parent.ChildCount) {
			utils.HandleError(c, http.StatusBadRequest, "Child bag limit reached for the parent bag")
			return
		}
		// Link child to parent
		bag.ParentBagID = req.ParentBagID
		bag.Linked = true // mark child as linked to its parent
		// Child's ChildCount remains 0 (not used for child bags)
	}

	// Save the new bag in the database
	if err := models.DB.Create(&bag).Error; err != nil {
		// Handle duplicate QR code or other DB errors
		log.Printf("Error registering bag (QRCode=%s): %v", req.QRCode, err)
		if isUniqueConstraintError(err) {
			utils.HandleError(c, http.StatusBadRequest, "A bag with this QR code already exists")
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "Could not register the bag")
		}
		return
	}

	// Success: return the created bag's details (e.g., ID and QRCode, BagType)
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Bag registered successfully",
		"bagId":       bag.ID,
		"qrCode":      bag.QRCode,
		"bagType":     bag.BagType,
		"childCount":  bag.ChildCount,
		"parentBagID": bag.ParentBagID,
	})
}

// isUniqueConstraintError checks if an error is caused by a unique constraint violation (e.g., duplicate QRCode).
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	// Check for common substrings in unique constraint errors (Postgres, SQLite, etc.)
	return (strpos(errMsg, "unique") || strpos(errMsg, "UNIQUE") || strpos(errMsg, "Duplicate entry"))
}

// strpos is a helper to check substring presence (case-sensitive)
func strpos(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && (stringIndex(s, substr) >= 0)
}

// stringIndex finds the index of substr in s or returns -1 if not found.
func stringIndex(s, substr string) int {
	// Simple implementation (could use strings.Index from standard library)
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func ListBags(c *gin.Context) {
	cursorStr := c.Query("cursor") // ID after which to start
	limitStr := c.Query("limit")   // max number of records to return
	var cursorID int
	var limit int

	if cursorStr != "" {
		if val, err := strconv.Atoi(cursorStr); err == nil {
			cursorID = val
		}
	}
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}
	if limit <= 0 || limit > 100 {
		limit = 10 // default limit (or enforce max limit)
	}

	var bags []models.Bag
	query := models.DB.Order("id ASC").Limit(limit)
	if cursorID > 0 {
		query = query.Where("id > ?", cursorID)
	}
	if err := query.Find(&bags).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Error listing bags")
		return
	}

	// Optionally, prepare next cursor (last ID from this page)
	var nextCursor *uint
	if len(bags) == limit {
		nextCursor = new(uint)
		*nextCursor = bags[len(bags)-1].ID
	}

	c.JSON(http.StatusOK, gin.H{
		"bags":       bags,
		"nextCursor": nextCursor, // clients can use this ID to fetch the next page
	})
}

func GetBagByQRCode(c *gin.Context) {
	code := c.Param("qr_code") // QR code from URL path or query
	var bag models.Bag
	if err := models.DB.Where("qr_code = ?", code).First(&bag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, http.StatusNotFound, "Bag not found")
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "Error retrieving bag details")
		}
		return
	}

	c.JSON(http.StatusOK, bag)
}
