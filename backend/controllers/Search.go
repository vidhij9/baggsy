package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"baggsy/backend/models"
	"baggsy/backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchBag handles looking up a bag (by QR code) to find its associated Bill (if any).
// It expects a QR code as a URL parameter.
func SearchBag(c *gin.Context) {
	qrCode := c.Param("qrcode")
	if qrCode == "" {
		utils.HandleError(c, http.StatusBadRequest, "QR code must be provided")
		return
	}

	var bag models.Bag
	// Look up the bag by its QR code
	if err := models.DB.Where("qr_code = ?", qrCode).First(&bag).Error; err != nil {
		log.Printf("Search failed for QRCode=%s: %v", qrCode, err)
		utils.HandleError(c, http.StatusNotFound, "Bag not found")
		return
	}

	// Determine action based on bag type
	if bag.BagType == "Child" {
		// For a child bag, find the parent and get its bill info
		if !bag.Linked || bag.ParentBagID == nil {
			// Child bag is not linked to any parent (should not normally happen if properly registered)
			utils.HandleError(c, http.StatusBadRequest, "This child bag is not linked to a parent")
			return
		}
		var parent models.Bag
		if err := models.DB.First(&parent, *bag.ParentBagID).Error; err != nil {
			log.Printf("Parent lookup failed for child bag ID=%d: %v", bag.ID, err)
			utils.HandleError(c, http.StatusInternalServerError, "Parent bag information not found")
			return
		}
		if !parent.LinkedToBill || parent.BillID == nil {
			// Parent exists but is not linked to any bill yet
			utils.HandleError(c, http.StatusOK, "Parent bag is not linked to a bill")
			return
		}
		// Return the bill ID associated with the parent bag
		c.JSON(http.StatusOK, gin.H{
			"message": "Bill found for this bag",
			"billId":  *parent.BillID,
		})
	} else if bag.BagType == "Parent" {
		// For a parent bag, directly check its bill linkage
		if !bag.LinkedToBill || bag.BillID == nil {
			utils.HandleError(c, http.StatusOK, "This parent bag is not linked to any bill")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Bill found for this bag",
				"billId":  *bag.BillID,
			})
		}
	} else {
		// In case BagType is somehow not set correctly
		utils.HandleError(c, http.StatusBadRequest, "Invalid bag type")
	}
}

// SearchBill handles GET /bills/{id}
func SearchBill(c *gin.Context) {
	idStr := c.Param("id") // Bill ID from URL param (as string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid bill ID")
		return
	}

	var bill models.Bill
	if err := models.DB.Where("id = ?", id).First(&bill).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, http.StatusNotFound, "Bill not found")
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "Error retrieving bill")
		}
		return
	}

	// Successfully found the bill
	c.JSON(http.StatusOK, bill)
}
