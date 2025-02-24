package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.AutoMigrate(&models.Bag{})
	return db
}

func TestListBagsHandler_EmptyDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.DB = setupTestDB()
	defer db.DB.Close()

	router := gin.Default()
	router.GET("/api/bags", ListBagsHandler)

	req, _ := http.NewRequest("GET", "/api/bags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "[]" {
		t.Errorf("Expected empty array '[]', got %s", w.Body.String())
	}
	if w.Header().Get("X-Total-Count") != "0" {
		t.Errorf("Expected X-Total-Count 0, got %s", w.Header().Get("X-Total-Count"))
	}
}

func TestListBagsHandler_WithBags(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.DB = setupTestDB()
	defer db.DB.Close()

	db.DB.Create(&models.Bag{QRCode: "P123-2", Type: "parent", ChildCount: 2})

	router := gin.Default()
	router.GET("/api/bags", ListBagsHandler)

	req, _ := http.NewRequest("GET", "/api/bags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Header().Get("X-Total-Count") != "1" {
		t.Errorf("Expected X-Total-Count 1, got %s", w.Header().Get("X-Total-Count"))
	}
}
