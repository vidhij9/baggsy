package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var redisClient = redis.NewClient(&redis.Options{Addr: "redis:6379"})

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

	// Validate QR code format (e.g., "P123-10")
	parts := strings.Split(parent.QRCode, "-")
	if len(parts) != 2 || !strings.HasPrefix(parts[0], "P") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid QR format. Use P<Number>-<ChildCount>"})
		return
	}
	count, err := strconv.Atoi(parts[1])
	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent bag must have at least one child"})
		return
	}
	parent.ChildCount = count

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
	if child.Type != "child" || child.ParentID == nil || child.QRCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must specify child type, parent ID, and QR code"})
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
	endDate := c.Query("endDate") // Fixed typo: was "startDate"
	unlinked := c.Query("unlinked") == "true"
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := db.DB.Model(&models.Bag{})
	if bagType != "" && (bagType == "parent" || bagType == "child") {
		query = query.Where("type = ?", bagType)
	}
	if startDate != "" && endDate != "" {
		if _, err := time.Parse("2006-01-02", startDate); err == nil {
			if _, err := time.Parse("2006-01-02", endDate); err == nil {
				query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
			}
		}
	}
	if unlinked {
		query = query.Where("linked = false AND type = 'parent'")
	}

	var total int64
	query.Count(&total)

	var bags []models.Bag
	if err := query.Offset(offset).Limit(limit).Find(&bags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bags"})
		return
	}

	// Manually fetch children for parent bags
	for i := range bags {
		if bags[i].Type == "parent" {
			var children []models.Bag
			db.DB.Model(&models.Bag{}).Where("parent_id = ?", bags[i].ID).Find(&children)
			bags[i].Children = children
		}
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
			resp.Children = bag.Children
			var link models.Link
			if db.DB.Where("parent_id = ?", bag.ID).First(&link).Error == nil {
				resp.BillID = link.BillID
			}
		} else if bag.ParentID != nil {
			var parent models.Bag
			if err := db.DB.First(&parent, *bag.ParentID).Error; err == nil {
				resp.ParentQR = parent.QRCode
			}
			var link models.Link
			if db.DB.Where("parent_id = ?", *bag.ParentID).First(&link).Error == nil {
				resp.BillID = link.BillID
			}
		}
		response = append(response, resp)
	}

	if len(response) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No bags found."})
	} else {
		c.Header("X-Total-Count", strconv.FormatInt(total, 10))
		c.JSON(http.StatusOK, response)
	}
}

func ListUnlinkedParentsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	db.DB.Model(&models.Bag{}).Where("type = 'parent' AND linked = false").Count(&total)

	var parents []models.Bag
	if err := db.DB.Where("type = 'parent' AND linked = false").Offset(offset).Limit(limit).Find(&parents).Error; err != nil {
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

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, response)
}

func FindChildBagsByParentQRHandler(c *gin.Context) {
	parentQR := c.Param("parentQR")
	var parent models.Bag
	if err := db.DB.Where("qr_code = ?", parentQR).First(&parent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent bag not found with QR: " + parentQR})
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
	if qr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code is required"})
		return
	}
	cacheKey := "bag:" + qr
	cached, err := redisClient.Get(cacheKey).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"data": cached})
		return
	}

	var bag models.Bag
	if err := db.DB.Where("lower(qr_code) = lower(?)", qr).First(&bag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bag not found with QR: " + qr})
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
		if err := db.DB.First(&parent, *bag.ParentID).Error; err == nil {
			response.ParentQR = parent.QRCode
		}
		var link models.Link
		if db.DB.Where("parent_id = ?", *bag.ParentID).First(&link).Error == nil {
			response.BillID = link.BillID
		}
	}

	redisClient.Set(cacheKey, bag, time.Hour)
	c.JSON(http.StatusOK, response)
}
