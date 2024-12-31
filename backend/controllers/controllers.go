package controllers

import (
	"baggsy/backend/database"
	"baggsy/backend/models"
	"errors"
	"fmt"
	"log"
	"net/http"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Example function to extract the number of child bags from the QR code
func extractChildBagCount(qrCode string) (int, error) {
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
func validateBagInput(qrCode, bagType string) error {
	if qrCode == "" || bagType == "" {
		return errors.New("QR Code and Bag Type are required")
	}
	if bagType != "Parent" && bagType != "Child" {
		return errors.New("bag Type must be 'Parent' or 'Child'")
	}
	return nil
}

// Centralized error handler
func handleError(c *gin.Context, statusCode int, message string, err error) {
	log.Printf("Error: %v", err)
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

// Register a bag and its child bags if applicable
func RegisterBag(c *gin.Context) {
	var bag models.Bag
	if err := c.BindJSON(&bag); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if err := validateBagInput(bag.QRCode, bag.BagType); err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if exists := database.DB.Where("qr_code = ?", bag.QRCode).First(&models.Bag{}).Error == nil; exists {
		handleError(c, http.StatusConflict, "Bag with this QR Code already exists", nil)
		return
	}

	childBagCount, err := extractChildBagCount(bag.QRCode)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	for i := 0; i < childBagCount; i++ {
		childBag := models.Bag{
			QRCode:  fmt.Sprintf("%s-Child-%d", bag.QRCode, i),
			BagType: "Child",
		}
		if err := database.DB.Create(&childBag).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to create child bags", err)
			return
		}
	}

	if err := database.DB.Create(&bag).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to register bag", err)
		return
	}

	log.Printf("Action: RegisterBag | QRCode: %s | BagType: %s", bag.QRCode, bag.BagType)
	c.JSON(http.StatusCreated, gin.H{"message": "Bag registered successfully"})
}

// Link parent bag and child bag, removing the child bag from the database
func LinkBags(c *gin.Context) {
	var link models.BagMap
	if err := c.BindJSON(&link); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if link.ParentBag == "" || link.ChildBag == "" {
		handleError(c, http.StatusBadRequest, "Parent Bag and Child Bag QR Codes are required", nil)
		return
	}

	// Validate child-parent relationship and max child bags
	if err := ValidateChildParentRelationship(link.ParentBag, link.ChildBag); err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := ValidateChildBagCount(link.ParentBag, 10); err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Link the bags
	if err := tx.Create(&link).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to link bags", err)
		tx.Rollback()
		return
	}

	// Soft delete the child bag
	if err := tx.Model(&models.Bag{}).Where("qr_code = ?", link.ChildBag).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to soft delete child bag", err)
		tx.Rollback()
		return
	}

	tx.Commit()
	log.Printf("Action: LinkBags | ParentBag: %s | ChildBag: %s", link.ParentBag, link.ChildBag)
	c.JSON(http.StatusOK, gin.H{"message": "Bags linked and child bag soft-deleted successfully"})
}

// Link parent bag to a bill and remove the parent bag from the database
func LinkBagToBill(c *gin.Context) {
	var link models.Link
	if err := c.BindJSON(&link); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if link.ParentBag == "" || link.BillID == "" {
		handleError(c, http.StatusBadRequest, "Parent Bag and Bill ID are required", nil)
		return
	}

	tx := database.DB.Begin()

	// Ensure the parent bag exists
	var parentBag models.Bag
	if err := tx.Where("qr_code = ? AND deleted_at IS NULL", link.ParentBag).First(&parentBag).Error; err != nil {
		handleError(c, http.StatusBadRequest, "Parent bag does not exist", err)
		tx.Rollback()
		return
	}

	// Link the parent bag to the bill
	if err := tx.Create(&link).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to link parent bag to bill", err)
		tx.Rollback()
		return
	}

	// Soft delete the parent bag
	if err := tx.Model(&models.Bag{}).Where("qr_code = ?", link.ParentBag).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to soft delete parent bag", err)
		tx.Rollback()
		return
	}

	tx.Commit()
	log.Printf("Action: LinkBagToBill | ParentBag: %s | BillID: %s", link.ParentBag, link.BillID)
	c.JSON(http.StatusOK, gin.H{"message": "Parent bag linked to bill and soft-deleted successfully"})
}

// Search for a bill by a child bag's QR Code
func SearchBillByBag(c *gin.Context) {
	qrCode := c.Query("qr_code")
	if qrCode == "" {
		handleError(c, http.StatusBadRequest, "QR Code is required", nil)
		return
	}

	var link models.Link
	if err := database.DB.Table("links").
		Joins("JOIN bag_map ON bag_map.parent_bag = links.parent_bag").
		Where("bag_map.child_bag = ? AND bag_map.deleted_at IS NULL", qrCode).
		First(&link).Error; err != nil {
		handleError(c, http.StatusNotFound, "Bill ID not found for this child bag", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"bill_id": link.BillID})
}

func UnlinkChildBag(c *gin.Context) {
	var link models.BagMap
	if err := c.BindJSON(&link); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	tx := database.DB.Begin()

	// Remove link
	if err := tx.Delete(&models.BagMap{}, "parent_bag = ? AND child_bag = ?", link.ParentBag, link.ChildBag).Error; err != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "Failed to unlink child bag", err)
		return
	}

	// Restore child bag
	restoredChild := models.Bag{QRCode: link.ChildBag, BagType: "Child"}
	if err := tx.Create(&restoredChild).Error; err != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "Failed to restore child bag", err)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Child bag unlinked and restored successfully"})
}
