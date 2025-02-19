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

// LinkBagToBillRequest represents the expected request body to link a bag to a bill.
type LinkBagToBillRequest struct {
	BillID uint   `json:"billId"`        // The bill ID to link to
	QRCode string `json:"parentBagCode"` // QR code of the parent bag to be linked
}

// LinkBagToBill links a parent bag to a given bill ID, marking it as used.
func LinkBagToBill(c *gin.Context) {
	var req LinkBagToBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Ensure the required fields are present
	if req.QRCode == "" || req.BillID == 0 {
		utils.HandleError(c, http.StatusBadRequest, "Both parent bag QR code and bill ID are required")
		return
	}

	// Find the bag by QR code
	var bag models.Bag
	if err := models.DB.Where("qr_code = ?", req.QRCode).First(&bag).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "Parent bag not found")
		return
	}

	// Only parent bags should be linked to bills
	if bag.BagType != "Parent" {
		utils.HandleError(c, http.StatusBadRequest, "Only parent bags can be linked to a bill")
		return
	}
	if bag.LinkedToBill {
		// Already linked to a bill, prevent duplicate linking
		utils.HandleError(c, http.StatusBadRequest, "This parent bag is already linked to a bill")
		return
	}

	// (Optional) Ensure all child bags are registered before linking to a bill.
	// This check is business-dependent; uncomment if needed.
	/*
	   var childCount int64
	   models.DB.Model(&models.Bag{}).Where("parent_bag_id = ?", bag.ID).Count(&childCount)
	   if childCount < int64(bag.ChildCount) {
	       utils.HandleError(c, http.StatusBadRequest, "Not all child bags have been registered for this parent")
	       return
	   }
	*/

	// Link the parent bag to the bill
	bag.BillID = &req.BillID
	bag.LinkedToBill = true
	if err := models.DB.Save(&bag).Error; err != nil {
		log.Printf("Failed to link bag (QRCode=%s) to BillID=%d: %v", req.QRCode, req.BillID, err)
		utils.HandleError(c, http.StatusInternalServerError, "Could not link bag to the bill")
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":       "Bag linked to bill successfully",
		"parentBagId":   bag.ID,
		"parentBagCode": bag.QRCode,
		"billId":        req.BillID,
	})
}

func UnlinkBagFromBill(c *gin.Context) {
	bagIdStr := c.Param("bag_id")
	bagID, err := strconv.Atoi(bagIdStr)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid bag ID")
		return
	}

	var bag models.Bag
	if err := models.DB.Where("id = ?", bagID).First(&bag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, http.StatusNotFound, "Bag not found")
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "Error finding bag")
		}
		return
	}

	// Set linked_to_bill = false instead of deleting the bag record
	if err := models.DB.Model(&bag).
		Updates(map[string]interface{}{"linked_to_bill": false}).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to unlink bag from bill")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bag unlinked from bill successfully"})
}

func GetParentBags(c *gin.Context) {
	billIdStr := c.Param("bill_id")
	billID, err := strconv.Atoi(billIdStr)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid bill ID")
		return
	}

	var parentBags []models.Bag
	if err := models.DB.Where("bill_id = ? AND linked_to_bill = ?", billID, true).
		Find(&parentBags).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Error retrieving parent bags")
		return
	}

	// Return the list of parent bags (empty list if none found)
	c.JSON(http.StatusOK, parentBags)
}
