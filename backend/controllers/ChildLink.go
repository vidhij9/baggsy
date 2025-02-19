package controllers

import (
	"baggsy/backend/models"
	"baggsy/backend/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LinkChildBagRequest defines the expected JSON body for linking child bags.
type LinkChildBagRequest struct {
	ParentBagID uint   `json:"parentBagId"` // ID of the parent bag
	ChildBagIDs []uint `json:"childBagIds"` // IDs of child bags to link
}

// LinkChildBag handles linking one or more child bags to a parent bag.
func LinkChildBag(c *gin.Context) {
	var req LinkChildBagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	parentID := req.ParentBagID
	childIDs := req.ChildBagIDs

	// Syntax Fix: proper if-condition outside of Printf for logging
	if parentID == 0 {
		log.Println("[LinkChildBag] Error: ParentBagID is missing or invalid (0).")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent bag ID is required"})
		return
	} else {
		log.Printf("[LinkChildBag] Linking child bags %v to ParentBagID %d\n", childIDs, parentID)
	}

	if len(childIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No child bag IDs provided to link"})
		return
	}

	// Begin a transaction for batch linking
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// Lock the parent bag row for update to prevent concurrent modifications
		var parent models.Bag
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&parent, parentID).Error; err != nil {
			return errors.New("parent bag not found")
		}

		// Validate parent bag
		if parent.BagType != "parent" {
			return errors.New("the specified bag is not a parent bag")
		}
		if parent.Linked {
			return errors.New("parent bag is already linked to a bill or closed for new children")
		}
		if parent.ChildCount < 1 {
			// ChildCount <= 0 means no children allowed (just a safety check)
			return errors.New("parent bag is not configured to have child bags")
		}

		// Count current number of linked children for this parent
		var currentChildren int64
		if err := tx.Model(&models.Bag{}).
			Where("parent_bag_id = ? AND linked = TRUE", parentID).
			Count(&currentChildren).Error; err != nil {
			return err // SQL error
		}
		// Ensure adding these children will not exceed capacity
		if currentChildren+int64(len(childIDs)) > int64(parent.ChildCount) {
			return errors.New("linking these child bags would exceed the parent bag's capacity")
		}

		// Loop through each child bag ID for linking
		for _, childID := range childIDs {
			// Lock and retrieve child bag record
			var child models.Bag
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&child, childID).Error; err != nil {
				return errors.New("child bag with ID " + fmt.Sprint(childID) + " not found")
			}

			// Validate child bag status
			if child.BagType != "child" {
				return errors.New("bag " + strconv.Itoa(int(child.ID)) + " is not a child-type bag")
				// Identifier() could be a method returning a human-readable identity (e.g., QR code or ID)
			}
			if child.ParentBagID != nil && *child.ParentBagID != 0 && child.Linked {
				// Child is currently linked to a parent
				return errors.New("child bag " + strconv.Itoa(int(child.ID)) + " is already linked to a parent bag")
			}

			// All good – link the child to the parent
			child.ParentBagID = &parentID // assign parent relationship
			child.Linked = true           // mark as linked
			if err := tx.Save(&child).Error; err != nil {
				return errors.New("failed to update child bag " + strconv.Itoa(int(child.ID)))
			}
		}

		// (Optional) If needed, update parent bag's record (e.g., mark registration complete)
		// Here we could update a field if one existed, but capacity enforcement via count is sufficient.

		return nil // commit transaction
	})

	if err != nil {
		// Transaction failed, respond with the error (rolled back automatically)
		log.Printf("[LinkChildBag] Linking failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// Success: all child bags linked
		c.JSON(http.StatusOK, gin.H{
			"message":          "Child bags linked successfully",
			"parentBagId":      parentID,
			"linkedChildCount": len(childIDs),
		})
	}
}

func UnlinkChildBag(c *gin.Context) {
	parentIdStr := c.Param("parent_id")
	childIdStr := c.Param("child_id")
	parentID, err := strconv.Atoi(parentIdStr)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid parent bag ID")
		return
	}
	childID, err := strconv.Atoi(childIdStr)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid child bag ID")
		return
	}

	// Find the child bag that is linked to the given parent
	var childBag models.Bag
	if err := models.DB.Where("id = ? AND parent_bag_id = ?", childID, parentID).
		First(&childBag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, http.StatusNotFound, "Child bag not found under the specified parent")
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "Error finding child bag")
		}
		return
	}

	// Set linked = false instead of removing the child-parent relationship
	if err := models.DB.Model(&childBag).
		Updates(map[string]interface{}{"linked": false}).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to unlink child bag from parent")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Child bag unlinked from parent successfully"})
}

func GetChildBags(c *gin.Context) {
	parentIdStr := c.Param("parent_id")
	parentID, err := strconv.Atoi(parentIdStr)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid parent bag ID")
		return
	}

	var childBags []models.Bag
	if err := models.DB.Where("parent_bag_id = ? AND linked = ?", parentID, true).
		Find(&childBags).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Error retrieving child bags")
		return
	}

	c.JSON(http.StatusOK, childBags)
}
