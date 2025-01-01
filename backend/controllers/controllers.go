package controllers

import (
	"baggsy/backend/database"
	"baggsy/backend/models"
	"baggsy/backend/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register a bag and its child bags if applicable
func RegisterBag(c *gin.Context) {
	var bag models.Bag
	if err := c.BindJSON(&bag); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	// Validate and check for duplicates
	if err := utils.ValidateBagInput(bag.QRCode, bag.BagType); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if exists := database.DB.Where("qr_code = ?", bag.QRCode).First(&models.Bag{}).Error == nil; exists {
		utils.HandleError(c, http.StatusConflict, "Bag with this QR Code already exists", nil)
		return
	}

	// Extract and batch-insert child bags
	childBagCount, err := utils.ExtractChildBagCount(bag.QRCode)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	childBags := make([]models.Bag, childBagCount)
	for i := 0; i < childBagCount; i++ {
		childBags[i] = models.Bag{
			QRCode:  fmt.Sprintf("%s-Child-%d", bag.QRCode, i),
			BagType: "Child",
		}
	}

	if err := database.DB.Create(&childBags).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to batch-insert child bags", err)
		return
	}

	// Register the parent bag
	if err := database.DB.Create(&bag).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to register bag", err)
		return
	}

	log.Printf("Action: RegisterBag | QRCode: %s | BagType: %s | ChildBagCount: %d", bag.QRCode, bag.BagType, childBagCount)
	c.JSON(http.StatusCreated, gin.H{"message": "Bag registered successfully", "child_bag_count": childBagCount})
}

// Link parent bag and child bag, removing the child bag from the database
func LinkBags(c *gin.Context) {
	var link models.BagMap
	if err := c.BindJSON(&link); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if link.ParentBag == "" || link.ChildBag == "" {
		utils.HandleError(c, http.StatusBadRequest, "Parent Bag and Child Bag QR Codes are required", nil)
		return
	}

	// Validate child-parent relationship and max child bags
	if err := utils.ValidateChildParentRelationship(link.ParentBag, link.ChildBag); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := utils.ValidateChildBagCount(link.ParentBag, 10); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
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
		utils.HandleError(c, http.StatusInternalServerError, "Failed to link bags", err)
		tx.Rollback()
		return
	}

	// Soft delete the child bag
	if err := tx.Model(&models.Bag{}).Where("qr_code = ?", link.ChildBag).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to soft delete child bag", err)
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
		utils.HandleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if link.ParentBag == "" || link.BillID == "" {
		utils.HandleError(c, http.StatusBadRequest, "Parent Bag and Bill ID are required", nil)
		return
	}

	tx := database.DB.Begin()

	// Ensure the parent bag exists
	var parentBag models.Bag
	if err := tx.Where("qr_code = ? AND deleted_at IS NULL", link.ParentBag).First(&parentBag).Error; err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Parent bag does not exist", err)
		tx.Rollback()
		return
	}

	// Link the parent bag to the bill
	if err := tx.Create(&link).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to link parent bag to bill", err)
		tx.Rollback()
		return
	}

	// Soft delete the parent bag
	if err := tx.Model(&models.Bag{}).Where("qr_code = ?", link.ParentBag).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to soft delete parent bag", err)
		tx.Rollback()
		return
	}

	tx.Commit()
	log.Printf("Action: LinkBagToBill | ParentBag: %s | BillID: %s", link.ParentBag, link.BillID)
	c.JSON(http.StatusOK, gin.H{"message": "Parent bag linked to bill and soft-deleted successfully"})
}

func SearchBillByBag(c *gin.Context) {
	qrCode := c.Query("qr_code")
	if qrCode == "" {
		utils.HandleError(c, http.StatusBadRequest, "QR Code is required", nil)
		return
	}

	var link models.Link
	if err := database.DB.Table("links").
		Joins("JOIN bag_map ON bag_map.parent_bag = links.parent_bag").
		Where("bag_map.child_bag = ? AND bag_map.deleted_at IS NULL", qrCode).
		First(&link).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "Bill ID not found for this child bag", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"bill_id": link.BillID})
}

func UnlinkChildBag(c *gin.Context) {
	var link models.BagMap
	if err := c.BindJSON(&link); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	tx := database.DB.Begin()

	// Remove link
	if err := tx.Delete(&models.BagMap{}, "parent_bag = ? AND child_bag = ?", link.ParentBag, link.ChildBag).Error; err != nil {
		tx.Rollback()
		utils.HandleError(c, http.StatusInternalServerError, "Failed to unlink child bag", err)
		return
	}

	// Restore child bag
	restoredChild := models.Bag{QRCode: link.ChildBag, BagType: "Child"}
	if err := tx.Create(&restoredChild).Error; err != nil {
		tx.Rollback()
		utils.HandleError(c, http.StatusInternalServerError, "Failed to restore child bag", err)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Child bag unlinked and restored successfully"})
}

func GetLinkedBagsByParent(c *gin.Context) {
	parentBag := c.Query("parent_bag")
	if parentBag == "" {
		utils.HandleError(c, http.StatusBadRequest, "Parent Bag QR Code is required", nil)
		return
	}

	var linkedBags []models.BagMap
	if err := database.DB.Where("parent_bag = ?", parentBag).Find(&linkedBags).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve linked bags", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": linkedBags})
}
