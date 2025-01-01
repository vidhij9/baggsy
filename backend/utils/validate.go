package utils

import (
	"baggsy/backend/database"
	"baggsy/backend/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Example function to extract the number of child bags from the QR code
func ExtractChildBagCount(qrCode string) (int, error) {
	// Assuming the QR code contains metadata in the format "parent-<childCount>"
	// Example: "parent-5" means 5 child bags
	parts := strings.Split(qrCode, "-")
	if len(parts) < 2 {
		return 0, errors.New("invalid QR code format")
	}

	childCount, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("invalid child count in QR code")
	}

	if childCount < 0 {
		return 0, errors.New("child count cannot be negative")
	}

	return childCount, nil
}

// Utility to validate bag inputs
func ValidateBagInput(qrCode, bagType string) error {
	if qrCode == "" || bagType == "" {
		return errors.New("QR Code and Bag Type are required")
	}
	if bagType != "Parent" && bagType != "Child" {
		return errors.New("bag Type must be 'Parent' or 'Child'")
	}
	return nil
}

// Centralized error handler
func HandleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	c.JSON(statusCode, gin.H{"error": message})
}

// Helper to validate child-parent relationships
func ValidateChildParentRelationship(parentBagQR string, childBagQR string) error {
	var bagMap models.BagMap
	if err := database.DB.Where("parent_bag = ? AND child_bag = ?", parentBagQR, childBagQR).First(&bagMap).Error; err == nil {
		return errors.New("child bag already linked to this parent bag")
	}
	return nil
}

// Helper to validate maximum child bag count
func ValidateChildBagCount(parentBag string, maxCount int) error {
	var count int64
	database.DB.Model(&models.BagMap{}).Where("parent_bag = ?", parentBag).Count(&count)
	if int(count) >= maxCount {
		return fmt.Errorf("parent bag already has the maximum number of child bags (%d)", maxCount)
	}
	return nil
}
