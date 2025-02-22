package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterParentHandler(c *gin.Context) {
	var parent models.Bag
	if err := c.ShouldBindJSON(&parent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if parent.Type != "parent" || parent.QRCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must specify parent type and QR code"})
		return
	}

	// Extract child count from QR code (e.g., "P123-10" means 10 children)
	parts := strings.Split(parent.QRCode, "-")
	if len(parts) > 1 {
		if count, err := strconv.Atoi(parts[1]); err == nil && count > 0 {
			parent.ChildCount = count
		}
	}

	tx := db.DB.Begin()
	if !tx.Where("qr_code = ?", parent.QRCode).First(&models.Bag{}).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "QR code already registered"})
		return
	}

	if err := tx.Create(&parent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register parent bag"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, parent)
}

func RegisterChildHandler(c *gin.Context) {
	var child models.Bag
	if err := c.ShouldBindJSON(&child); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if child.Type != "child" || child.ParentID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must specify child type and parent ID"})
		return
	}

	tx := db.DB.Begin()
	var parent models.Bag
	if err := tx.First(&parent, *child.ParentID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent not found"})
		return
	}

	var childCount int64
	tx.Model(&models.Bag{}).Where("parent_id = ?", *child.ParentID).Count(&childCount)
	if int(childCount) >= parent.ChildCount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent capacity reached"})
		return
	}

	if !tx.Where("qr_code = ?", child.QRCode).First(&models.Bag{}).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "QR code already registered"})
		return
	}

	if err := tx.Create(&child).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register child bag"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"child": child, "currentCount": childCount + 1, "capacity": parent.ChildCount})
}

func ListBagsHandler(c *gin.Context) {
	bagType := c.Query("type")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	unlinked := c.Query("unlinked") == "true"

	query := db.DB.Model(&models.Bag{})
	if bagType != "" {
		query = query.Where("type = ?", bagType)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	if unlinked {
		query = query.Where("linked = false AND type = 'parent'")
	}

	var bags []models.Bag
	if err := query.Find(&bags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bags"})
		return
	}

	type BagResponse struct {
		Bag      models.Bag   `json:"bag"`
		Children []models.Bag `json:"children,omitempty"`
		BillID   string       `json:"billID,omitempty"`
		ParentQR string       `json:"parentQR,omitempty"`
	}
	var response []BagResponse

	for _, bag := range bags {
		resp := BagResponse{Bag: bag}
		if bag.Type == "parent" {
			var children []models.Bag
			db.DB.Where("parent_id = ?", bag.ID).Find(&children)
			resp.Children = children
			var link models.Link
			if db.DB.Where("parent_id = ?", bag.ID).First(&link).Error == nil {
				resp.BillID = link.BillID
			}
		} else if bag.ParentID != nil {
			var parent models.Bag
			db.DB.First(&parent, *bag.ParentID)
			resp.ParentQR = parent.QRCode
			var link models.Link
			if db.DB.Where("parent_id = ?", *bag.ParentID).First(&link).Error == nil {
				resp.BillID = link.BillID
			}
		}
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, response)
}

func ListUnlinkedParentsHandler(c *gin.Context) {
	var parents []models.Bag
	if err := db.DB.Where("type = 'parent' AND linked = false").Find(&parents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list unlinked parents"})
		return
	}

	type ParentResponse struct {
		Bag      models.Bag   `json:"bag"`
		Children []models.Bag `json:"children"`
	}
	var response []ParentResponse

	for _, parent := range parents {
		var children []models.Bag
		db.DB.Where("parent_id = ?", parent.ID).Find(&children)
		response = append(response, ParentResponse{Bag: parent, Children: children})
	}

	c.JSON(http.StatusOK, response)
}

func FindChildBagsByParentQRHandler(c *gin.Context) {
	parentQR := c.Param("parentQR")
	var parent models.Bag
	if err := db.DB.Where("qr_code = ?", parentQR).First(&parent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent bag not found"})
		return
	}

	var children []models.Bag
	if err := db.DB.Where("parent_id = ?", parent.ID).Find(&children).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find child bags"})
		return
	}

	var billID string
	var link models.Link
	if db.DB.Where("parent_id = ?", parent.ID).First(&link).Error == nil {
		billID = link.BillID
	}

	c.JSON(http.StatusOK, gin.H{"parent": parent, "children": children, "billID": billID})
}

func SearchBagByQRHandler(c *gin.Context) {
	qr := c.Param("qr")
	var bag models.Bag
	if err := db.DB.Where("qr_code = ?", qr).First(&bag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bag not found"})
		return
	}

	response := struct {
		Bag      models.Bag   `json:"bag"`
		Children []models.Bag `json:"children,omitempty"`
		ParentQR string       `json:"parentQR,omitempty"`
		BillID   string       `json:"billID,omitempty"`
	}{Bag: bag}

	if bag.Type == "parent" {
		var children []models.Bag
		db.DB.Where("parent_id = ?", bag.ID).Find(&children)
		response.Children = children
		var link models.Link
		if db.DB.Where("parent_id = ?", bag.ID).First(&link).Error == nil {
			response.BillID = link.BillID
		}
	} else if bag.ParentID != nil {
		var parent models.Bag
		db.DB.First(&parent, *bag.ParentID)
		response.ParentQR = parent.QRCode
		var link models.Link
		if db.DB.Where("parent_id = ?", *bag.ParentID).First(&link).Error == nil {
			response.BillID = link.BillID
		}
	}

	c.JSON(http.StatusOK, response)
}
