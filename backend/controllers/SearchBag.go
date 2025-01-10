package controllers

import (
	"baggsy/backend/database"
	"baggsy/backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /search-bag?qrCode=XYZ
func SearchBag(c *gin.Context) {
	qrCode := c.Query("qrCode")
	if qrCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qrCode is required"})
		return
	}

	// 1) Unscoped query so we can find the bag even if it's soft-deleted
	var bag models.Bag
	if err := database.DB.Unscoped().
		Where("qr_code = ?", qrCode).
		First(&bag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bag not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error retrieving bag"})
		}
		return
	}

	// 2) If it's a parent, find Bill in your "links" table
	if bag.BagType == "Parent" {
		// e.g. look up links table
		var link models.Link
		if err := database.DB.Where("parent_bag = ?", bag.QRCode).First(&link).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{"message": "No bill linked to this parent bag"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error retrieving link"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"bagType": "Parent",
			"billId":  link.BillID,
		})
		return
	}

	// 3) If it's a child, we might retrieve the parent's Bill ID
	if bag.BagType == "Child" {
		if bag.ParentBag == "" {
			c.JSON(http.StatusOK, gin.H{"message": "Child has no parent assigned"})
			return
		}
		// find link for that parent
		var link models.Link
		if err := database.DB.Where("parent_bag = ?", bag.ParentBag).First(&link).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{"message": "No bill linked to this child bag"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error retrieving link for child's parent"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"bagType": "Child",
			"billId":  link.BillID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unknown bag type"})
}
