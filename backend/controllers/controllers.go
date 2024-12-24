package controllers

import (
	"baggsy/backend/database"
	"baggsy/backend/models"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var mu sync.Mutex // Mutex for thread-safe operations

func validateBagInput(qrCode, bagType string) error {
	if qrCode == "" || bagType == "" {
		return errors.New("QR Code and Bag Type are required")
	}
	if bagType != "Parent" && bagType != "Child" {
		return errors.New("bag Type must be 'Parent' or 'Child'")
	}
	return nil
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	log.Printf("Error: %v", err)
	c.JSON(statusCode, gin.H{"error": message})
}

func RegisterBag(c *gin.Context) {
	var body models.Bag
	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if err := validateBagInput(body.QRCode, body.BagType); err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := database.DB.Create(&body).Error
	if err != nil {
		handleError(c, http.StatusConflict, "Bag with this QR code already exists", err)
		return
	}

	log.Printf("Action: RegisterBag | QRCode: %s | BagType: %s", body.QRCode, body.BagType)
	c.JSON(http.StatusCreated, gin.H{"message": "Bag registered successfully"})
}

func LinkBags(c *gin.Context) {
	var body models.Link
	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if body.ParentBag == "" || body.ChildBag == "" {
		handleError(c, http.StatusBadRequest, "Parent Bag and Child Bag QR Codes are required", nil)
		return
	}

	if body.ParentBag == body.ChildBag {
		handleError(c, http.StatusBadRequest, "Parent and Child Bag cannot be the same", nil)
		return
	}

	var parentBag, childBag models.Bag
	if err := database.DB.Where("qr_code = ?", body.ParentBag).First(&parentBag).Error; err != nil {
		handleError(c, http.StatusBadRequest, "Parent bag does not exist", err)
		return
	}
	if err := database.DB.Where("qr_code = ?", body.ChildBag).First(&childBag).Error; err != nil {
		handleError(c, http.StatusBadRequest, "Child bag does not exist", err)
		return
	}

	err := database.DB.Create(&body).Error
	if err != nil {
		handleError(c, http.StatusConflict, "Link between parent and child already exists", err)
		return
	}

	log.Printf("Action: LinkBags | ParentBag: %s | ChildBag: %s", body.ParentBag, body.ChildBag)
	c.JSON(http.StatusOK, gin.H{"message": "Bags linked successfully"})
}

func LinkBagToBill(c *gin.Context) {
	var body models.BagMap
	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if body.ParentBag == "" || body.BillID == "" {
		handleError(c, http.StatusBadRequest, "Parent Bag and Bill ID are required", nil)
		return
	}

	err := database.DB.Create(&body).Error
	if err != nil {
		handleError(c, http.StatusConflict, "Parent bag is already linked to a bill", err)
		return
	}

	log.Printf("Action: LinkBagToBill | ParentBag: %s | BillID: %s", body.ParentBag, body.BillID)
	c.JSON(http.StatusOK, gin.H{"message": "Bag linked to bill successfully"})
}

func EditBillID(c *gin.Context) {
	var body models.BagMap
	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if body.ParentBag == "" || body.BillID == "" {
		handleError(c, http.StatusBadRequest, "Parent Bag and new Bill ID are required", nil)
		return
	}

	err := database.DB.Model(&models.BagMap{}).Where("parent_bag = ?", body.ParentBag).Update("bill_id", body.BillID).Error
	if err != nil {
		handleError(c, http.StatusNotFound, "Parent bag not found", err)
		return
	}

	log.Printf("Action: EditBillID | ParentBag: %s | NewBillID: %s", body.ParentBag, body.BillID)
	c.JSON(http.StatusOK, gin.H{"message": "Bill ID updated successfully"})
}

func SearchBillByBag(c *gin.Context) {
	qrCode := c.Query("qr_code")
	if qrCode == "" {
		handleError(c, http.StatusBadRequest, "QR Code is required", nil)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var bagMap models.BagMap
	err := database.DB.Where("parent_bag = ?", qrCode).First(&bagMap).Error
	if err != nil {
		handleError(c, http.StatusNotFound, "Bag not linked to any bill", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"bill_id": bagMap.BillID})
}

func SearchBagsByBill(c *gin.Context) {
	billID := c.Query("bill_id")
	if billID == "" {
		handleError(c, http.StatusBadRequest, "Bill ID is required", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	mu.Lock()
	defer mu.Unlock()

	var linkedBags []models.BagMap
	err := database.DB.Where("bill_id = ?", billID).Limit(limit).Offset(offset).Find(&linkedBags).Error
	if err != nil {
		handleError(c, http.StatusNotFound, "No linked bags found for this Bill ID", err)
		return
	}

	var bagQRs []string
	for _, bag := range linkedBags {
		bagQRs = append(bagQRs, bag.ParentBag)
	}

	c.JSON(http.StatusOK, gin.H{"bags": bagQRs})
}
