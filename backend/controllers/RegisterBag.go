package controllers

import (
	"baggsy/backend/database"
	"baggsy/backend/models"
	"baggsy/backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func RegisterBag(c *gin.Context) {
	var bag models.Bag
	var err error

	// Parse the JSON input
	if err := c.ShouldBindJSON(&bag); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	// Validate input and check for duplicates
	if err := utils.ValidateBagInput(bag.QRCode, bag.BagType, bag.ChildCount); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if exists := database.DB.Where("qr_code = ?", bag.QRCode).First(&models.Bag{}).Error == nil; exists {
		utils.HandleError(c, http.StatusConflict, "Bag with this QR Code already exists", nil)
		return
	}

	// Handle Parent Bag Registration
	if bag.BagType == "Parent" {
		// Extract child count from QR code
		bag.ChildCount, err = utils.ExtractChildBagCount(bag.QRCode)
		if err != nil || bag.ChildCount <= 0 {
			utils.HandleError(c, http.StatusBadRequest, "Invalid child count in QR Code", nil)
			return
		}

		// Save the parent bag
		if err := database.DB.Create(&bag).Error; err != nil {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to register parent bag", err)
			return
		}

		log.Printf("Action: RegisterBag | ParentBag: %s | ChildCount: %d", bag.QRCode, bag.ChildCount)
		c.JSON(http.StatusCreated, gin.H{
			"message": "Parent bag registered successfully",
			"bag":     bag,
		})
		return
	}

	// Handle Child Bag Registration
	if bag.BagType == "Child" {
		// Ensure parent bag QR code is provided
		if bag.ParentBag == "" {
			utils.HandleError(c, http.StatusBadRequest, "Child bag must have a parent bag", nil)
			return
		}

		// Validate the parent bag
		var parentBag models.Bag
		if err := database.DB.Where("qr_code = ? AND bag_type = 'Parent'", bag.ParentBag).First(&parentBag).Error; err != nil {
			utils.HandleError(c, http.StatusNotFound, "Parent bag not found", err)
			return
		}

		// Save the child bag
		if err := database.DB.Create(&bag).Error; err != nil {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to register child bag", err)
			return
		}

		log.Printf("Action: RegisterBag | ParentBag: %s | ChildBag: %s", bag.ParentBag, bag.QRCode)
		c.JSON(http.StatusCreated, gin.H{
			"message": "Child bag registered successfully",
			"childBag": gin.H{
				"qrCode":    bag.QRCode,
				"parentBag": bag.ParentBag,
			},
		})
		return
	}

	// If BagType is invalid
	utils.HandleError(c, http.StatusBadRequest, "Invalid Bag Type", nil)
}

// func LinkChildBagToParent(c *gin.Context) {
// 	var bagRequest models.BagRequest
// 	var bag models.Bag

// 	// Parse and validate JSON payload
// 	if err := c.ShouldBindJSON(&bagRequest); err != nil {
// 		utils.HandleError(c, http.StatusBadRequest, "Invalid input", err)
// 		return
// 	}

// 	log.Printf("Parsed Request: ParentBag=%s, ChildBag=%s", bagRequest.ParentBag, bagRequest.ChildBag)

// 	// Ensure the parent bag exists
// 	var parentBag models.Bag
// 	if err := database.DB.Where("qr_code = ? AND bag_type = 'Parent'", bagRequest.ParentBag).First(&parentBag).Error; err != nil {
// 		utils.HandleError(c, http.StatusNotFound, "Parent bag not found", err)
// 		return
// 	}

// 	// Validate child count
// 	var linkedChildCount int64
// 	database.DB.Model(bag).Where("parent_bag = ?", parentBag.QRCode).Count(&linkedChildCount)

// 	if int(linkedChildCount) >= parentBag.ChildCount {
// 		utils.HandleError(c, http.StatusBadRequest, "Child bag limit exceeded for this parent bag", nil)
// 		return
// 	}

// 	// Save the child bag in the database
// 	childBag := models.Bag{
// 		QRCode:    bagRequest.ChildBag,
// 		BagType:   "Child",
// 		ParentBag: parentBag.QRCode,
// 	}
// 	if err := database.DB.Create(&childBag).Error; err != nil {
// 		utils.HandleError(c, http.StatusInternalServerError, "Failed to register child bag", err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message":          "Child bag linked to parent bag successfully",
// 		"linkedChildCount": linkedChildCount + 1,
// 	})
// }

// LinkBagToBill links a parent bag to a bill, ensuring concurrency safety via row locks,
// validating bag type, marking it as linked (and optionally soft-deleting), and providing
// detailed error messages.
func LinkBagToBilll(c *gin.Context) {
	var link models.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid JSON payload", err)
		return
	}

	// Validate request fields
	if link.ParentBag == "" || link.BillID == "" {
		utils.HandleError(c, http.StatusBadRequest, "Parent Bag and Bill ID are required", nil)
		return
	}

	// Begin transaction
	tx := database.DB.Begin()
	defer func() {
		// In case of panic or unexpected error
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 1) Find and lock the parent bag row FOR UPDATE to avoid concurrency issues.
	//    This will block other transactions attempting to update this same row
	//    until this transaction completes/rolls back.
	var parentBag models.Bag
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}). // row-level lock
		Where("qr_code = ? AND deleted_at IS NULL", link.ParentBag).
		First(&parentBag).Error; err != nil {

		// Differentiate between not found vs. other DB errors
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(c, http.StatusNotFound,
				"Parent bag not found or already deleted", err)
		} else {
			utils.HandleError(c, http.StatusInternalServerError,
				"Database error while retrieving parent bag", err)
		}
		tx.Rollback()
		return
	}

	// 2) Check Bag Type: Ensure it's a "Parent" bag.
	if parentBag.BagType != "Parent" {
		utils.HandleError(c, http.StatusBadRequest,
			"The specified bag is not a parent bag", nil)
		tx.Rollback()
		return
	}

	// Check if already linked
	if parentBag.Linked {
		utils.HandleError(c, http.StatusBadRequest,
			"This parent bag is already linked to a bill", nil)
		tx.Rollback()
		return
	}

	// 3) Create a record in the "links" table if your schema requires storing link details.
	//    (We're assuming 'models.Link' references ParentBag and BillID.)
	if err := tx.Create(&link).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError,
			"Failed to create link record for parent bag", err)
		tx.Rollback()
		return
	}

	// Mark the parent bag as linked. If your Bag model also has a BillID field,
	// you can set: parentBag.BillID = link.BillID
	parentBag.Linked = true
	// parentBag.BillID = link.BillID

	if err := tx.Save(&parentBag).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError,
			"Failed to mark the parent bag as linked", err)
		tx.Rollback()
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError,
			"Failed to commit the transaction", err)
		return
	}

	log.Printf("Action: LinkBagToBill | ParentBag: %s | BillID: %s",
		link.ParentBag, link.BillID)
	c.JSON(http.StatusOK, gin.H{"message": "Parent bag linked to bill successfully"})
}
