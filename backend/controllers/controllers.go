package controllers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"baggsy/models"

	"github.com/gin-gonic/gin"
)

var (
	bags         = []models.Bag{}
	bills        = []models.Bill{}
	distributors = []models.Distributor{}
	mutex        sync.Mutex // Mutex to prevent race conditions
)

// CreateBag handles adding a new bag and optionally linking it to a parent bag
func CreateBag(c *gin.Context) {
	var newBag models.Bag
	if err := c.ShouldBindJSON(&newBag); err != nil {
		log.Printf("Error parsing input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Check for unique QR Code
	for _, bag := range bags {
		if bag.QRCode == newBag.QRCode {
			log.Printf("Duplicate QR Code found: %v", newBag.QRCode)
			c.JSON(http.StatusConflict, gin.H{"error": "Bag with this QR Code already exists"})
			return
		}
	}

	if newBag.ParentBagID != nil {
		// Validate parent bag exists
		found := false
		for _, bag := range bags {
			if bag.ID == *newBag.ParentBagID {
				found = true
				break
			}
		}
		if !found {
			log.Printf("Parent bag ID not found: %v", *newBag.ParentBagID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent bag ID not found"})
			return
		}
	}

	bags = append(bags, newBag)
	log.Printf("Bag created successfully: %v", newBag)
	c.JSON(http.StatusCreated, gin.H{"message": "Bag created successfully", "bag": newBag})
}

// LinkBagsToBill links parent bags to an existing SAP bill ID
func LinkBagsToSAPBill(c *gin.Context) {
	var linkRequest struct {
		SAPBillID  string   `json:"sap_bill_id"`
		ParentBags []string `json:"parent_bags"`
	}
	if err := c.ShouldBindJSON(&linkRequest); err != nil {
		log.Printf("Error parsing input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Validate parent bags exist and are not linked to another SAP bill ID
	for _, parentID := range linkRequest.ParentBags {
		found := false
		for i, bag := range bags {
			if bag.ID == parentID && bag.BagType == "Parent" {
				if bag.BillID != nil {
					log.Printf("Parent bag %v is already linked to another SAP bill", parentID)
					c.JSON(http.StatusConflict, gin.H{"error": "Parent bag already linked to another SAP bill", "parent_bag_id": parentID})
					return
				}
				found = true
				bags[i].BillID = &linkRequest.SAPBillID // Link bag to the SAP bill
				break
			}
		}
		if !found {
			log.Printf("Parent bag ID not found or invalid: %v", parentID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent bag ID not found or invalid", "parent_bag_id": parentID})
			return
		}
	}

	log.Printf("Bags linked to SAP bill successfully: %v", linkRequest)
	c.JSON(http.StatusOK, gin.H{"message": "Bags linked to SAP bill successfully", "sap_bill_id": linkRequest.SAPBillID})
}

// GetBagsWithPagination retrieves bags with pagination
func GetBagsWithPagination(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	limitQuery := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	start := (page - 1) * limit
	end := start + limit

	mutex.Lock()
	total := len(bags)
	if start >= total {
		mutex.Unlock()
		c.JSON(http.StatusOK, gin.H{"bags": []models.Bag{}, "page": page, "limit": limit, "total": total})
		return
	}

	if end > total {
		end = total
	}

	paginatedBags := bags[start:end]
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"bags": paginatedBags, "page": page, "limit": limit, "total": total})
}
