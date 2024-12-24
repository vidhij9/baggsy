package controllers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	bags   = []map[string]string{}
	links  = []map[string]string{}
	bagMap = make(map[string]string) // bag -> bill
	mu     sync.Mutex                // Mutex for thread-safe operations
)

// Register a bag
func RegisterBag(c *gin.Context) {
	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	qrCode := body["qr_code"]
	bagType := body["bag_type"]

	mu.Lock()
	defer mu.Unlock()
	for _, bag := range bags {
		if bag["qr_code"] == qrCode {
			c.JSON(http.StatusConflict, gin.H{"error": "Bag with this QR code already exists"})
			return
		}
	}
	bags = append(bags, map[string]string{
		"qr_code":  qrCode,
		"bag_type": bagType,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "Bag registered successfully"})
}

// Link parent and child bags
func LinkBags(c *gin.Context) {
	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	parentBag := body["parent_bag_qr_code"]
	childBag := body["child_bag_qr_code"]

	mu.Lock()
	defer mu.Unlock()
	foundParent, foundChild := false, false
	for _, bag := range bags {
		if bag["qr_code"] == parentBag {
			foundParent = true
		}
		if bag["qr_code"] == childBag {
			foundChild = true
		}
	}
	if !foundParent || !foundChild {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either parent or child bag does not exist"})
		return
	}
	for _, link := range links {
		if link["parent"] == parentBag && link["child"] == childBag {
			c.JSON(http.StatusConflict, gin.H{"error": "Link between parent and child already exists"})
			return
		}
	}
	links = append(links, map[string]string{
		"parent": parentBag,
		"child":  childBag,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Bags linked successfully"})
}

// Link parent bag to a bill ID
func LinkBagToBill(c *gin.Context) {
	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	parentBag := body["parent_bag_qr_code"]
	billID := body["bill_id"]

	mu.Lock()
	defer mu.Unlock()
	if _, exists := bagMap[parentBag]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Parent bag is already linked to a bill"})
		return
	}
	bagMap[parentBag] = billID

	c.JSON(http.StatusOK, gin.H{"message": "Bag linked to bill successfully"})
}

// Edit bill ID for a parent bag
func EditBillID(c *gin.Context) {
	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	parentBag := body["parent_bag_qr_code"]
	newBillID := body["new_bill_id"]

	mu.Lock()
	defer mu.Unlock()
	if _, exists := bagMap[parentBag]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent bag not found"})
		return
	}
	bagMap[parentBag] = newBillID

	c.JSON(http.StatusOK, gin.H{"message": "Bill ID updated successfully"})
}

// Search for a bag to get its bill ID
func SearchBillByBag(c *gin.Context) {
	qrCode := c.Query("qr_code")
	mu.Lock()
	billID, exists := bagMap[qrCode]
	mu.Unlock()
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bag not linked to any bill"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bill_id": billID})
}

// Search for a bill ID to get all linked bags
func SearchBagsByBill(c *gin.Context) {
	billID := c.Query("bill_id")
	var linkedBags []string

	mu.Lock()
	for bag, bill := range bagMap {
		if bill == billID {
			linkedBags = append(linkedBags, bag)
		}
	}
	mu.Unlock()

	if len(linkedBags) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No bags found for this bill ID", "bill_id": billID})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bags": linkedBags})
}
