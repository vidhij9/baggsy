package controllers

import (
	"baggsy/backend/database"
	"baggsy/backend/models"
	"baggsy/backend/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BillRequest struct {
	ParentBag string `gorm:"not null" json:"parentBag" binding:"linkuired"`
	BillID    string `gorm:"not null" json:"billID" binding:"linkuired"`
}

func LinkBagToBill(c *gin.Context) {
	var link models.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Validate linkuest fields
	if link.ParentBag == "" || link.BillID == "" {
		utils.HandleError(c, http.StatusBadRequest, "Parent Bag and Bill ID are linkuired", nil)
		return
	}

	// 1) Unscoped query to see if the parent bag is already soft-deleted
	var parentBag models.Bag
	if err := database.DB.Unscoped().
		Where("qr_code = ? AND bag_type = 'Parent'", link.ParentBag).
		First(&parentBag).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent bag not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error retrieving parent bag"})
		}
		return
	}

	// 2) If it's soft deleted => already linked
	if parentBag.DeletedAt.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This parent bag is already linked to a bill"})
		return
	}

	// 3) Link the parent bag to Bill ID in your "links" table or however you store it
	newLink := models.Link{
		ParentBag: link.ParentBag,
		BillID:    link.BillID,
	}
	if err := database.DB.Create(&newLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link record"})
		return
	}

	// optional: set parentBag.Linked = true if you use that field
	parentBag.Linked = true
	if err := database.DB.Save(&parentBag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not mark parent bag as linked"})
		return
	}

	// 4) Soft-delete the parent bag so it can't be linked again
	if err := database.DB.Delete(&parentBag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft-delete parent bag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Parent bag linked to bill and soft-deleted successfully",
	})
}
