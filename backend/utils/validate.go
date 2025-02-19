package utils

import (
	"baggsy/backend/models"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExtractChildBagCount returns how many child bags are linked to the given parent bag.
func ExtractChildBagCount(db *gorm.DB, parentID uint) int {
	var count int64
	// Count all bags that have ParentBagID = parentID (i.e., are children of this parent)
	err := db.Model(&models.Bag{}).Where("parent_bag_id = ?", parentID).Count(&count).Error
	if err != nil {
		log.Println("ExtractChildBagCount: error counting child bags for parent", parentID, "-", err)
		return 0 // on error, return 0 (and log it); alternatively, handle error upstream
	}
	return int(count)
}

// Utility to validate bag inputs
func ValidateBagInput(qrCode, bagType string, childCount int) error {
	if qrCode == "" || bagType == "" {
		return errors.New("QR Code and Bag Type are required")
	}
	if bagType != "Parent" && bagType != "Child" {
		return errors.New("bag Type must be 'Parent' or 'Child'")
	}
	return nil
}

// Centralized error handler
func HandleError(c *gin.Context, statusCode int, errMsg string) {
	// Optionally log the error message here as well
	c.JSON(statusCode, gin.H{"error": errMsg})
}
