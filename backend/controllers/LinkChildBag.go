package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"baggsy/backend/database"
	"baggsy/backend/models"
)

type BagRequest struct {
	ParentBag string `json:"parentBag" binding:"required"` // Ensure this matches the JSON payload
	ChildBag  string `json:"childBag" binding:"required"`  // Ensure this matches the JSON payload
}

// LinkChildBag links a child bag to a parent bag, then soft-deletes the child bag.
func LinkChildBag(c *gin.Context) {
	var req BagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// 1) Check if the child bag is already soft-deleted (meaning it's linked already).
	//    We do an unscoped query so we see *all* rows, including soft-deleted.
	var childBag models.Bag
	if err := database.DB.Unscoped().
		Where("qr_code = ?", req.ChildBag).
		First(&childBag).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Child bag not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error finding child bag"})
		}
		return
	}

	// 2) If the child bag is soft-deleted, it means it's already linked.
	if childBag.DeletedAt.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This child bag is already linked to a parent bag"})
		return
	}

	// 3) If it's not soft-deleted but bagType isn't Child, fail:
	if childBag.BagType != "Child" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provided bag is not a child bag"})
		return
	}

	// 4) Check the parent bag (normal query, ignoring soft-deleted parent if needed).
	var parentBag models.Bag
	if err := database.DB.
		Where("qr_code = ? AND bag_type = 'Parent' AND deleted_at IS NULL", req.ParentBag).
		First(&parentBag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent bag not found or is soft-deleted"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error finding parent bag"})
		}
		return
	}

	// 5) Link child → parent
	//    e.g., childBag.ParentBag = parentBag.QRCode
	childBag.ParentBag = parentBag.QRCode
	if err := database.DB.Save(&childBag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update child bag's parent field"})
		return
	}

	// 6) Soft-delete the child bag so it can't be linked again.
	if err := database.DB.Delete(&childBag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft-delete child bag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Child bag successfully linked and soft-deleted",
	})
}
