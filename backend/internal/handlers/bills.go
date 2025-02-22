package handlers

import (
	"net/http"
	"strconv"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/models"

	"github.com/gin-gonic/gin"
)

func LinkBagsToBillHandler(c *gin.Context) {
	var req struct {
		BillID    string `json:"billID" binding:"required"`
		ParentIDs []uint `json:"parentIDs" binding:"required"`
		Capacity  int    `json:"capacity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(req.ParentIDs) != req.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Number of parent bags must match capacity"})
		return
	}

	tx := db.DB.Begin()
	for _, parentID := range req.ParentIDs {
		var bag models.Bag
		if err := tx.First(&bag, parentID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent bag not found"})
			return
		}
		if bag.Linked {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Bag already linked"})
			return
		}

		link := models.Link{ParentID: parentID, BillID: req.BillID}
		if err := tx.Create(&link).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link bag"})
			return
		}
		if err := tx.Model(&bag).Update("linked", true).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bag status"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Bags linked successfully"})
}

func UnlinkBagHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bag ID"})
		return
	}

	tx := db.DB.Begin()
	var link models.Link
	if err := tx.Where("parent_id = ?", id).First(&link).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	if err := tx.Delete(&link).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink bag"})
		return
	}

	if err := tx.Model(&models.Bag{}).Where("id = ?", id).Update("linked", false).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bag status"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Bag unlinked successfully"})
}

func ListBillsHandler(c *gin.Context) {
	var links []models.Link
	if err := db.DB.Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bills"})
		return
	}

	billMap := make(map[string][]models.Bag)
	for _, link := range links {
		var bag models.Bag
		db.DB.First(&bag, link.ParentID)
		billMap[link.BillID] = append(billMap[link.BillID], bag)
	}

	type BillResponse struct {
		BillID string       `json:"billID"`
		Bags   []models.Bag `json:"bags"`
	}
	var response []BillResponse
	for billID, bags := range billMap {
		response = append(response, BillResponse{BillID: billID, Bags: bags})
	}

	c.JSON(http.StatusOK, response)
}

func SearchBillByNumberHandler(c *gin.Context) {
	billID := c.Param("billID")
	var links []models.Link
	if err := db.DB.Where("bill_id = ?", billID).Find(&links).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	var bags []models.Bag
	for _, link := range links {
		var bag models.Bag
		if err := db.DB.First(&bag, link.ParentID).Error; err == nil {
			var children []models.Bag
			db.DB.Where("parent_id = ?", bag.ID).Find(&children)
			bag.Children = children // Assuming Bag model has a Children field for response
			bags = append(bags, bag)
		}
	}

	c.JSON(http.StatusOK, gin.H{"billID": billID, "bags": bags})
}
