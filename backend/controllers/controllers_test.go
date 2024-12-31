package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/register-bag", RegisterBag)
	router.POST("/link-bags", LinkBags)
	router.POST("/link-bag-to-bill", LinkBagToBill)
	router.GET("/search-bill-by-bag", SearchBillByBag)
	return router
}

func TestRegisterBag(t *testing.T) {
	router := setupRouter()

	body := map[string]string{
		"qr_code":  "12345",
		"bag_type": "Parent",
		"status":   "Active",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create-bag", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Bag registered successfully")
}

func TestDuplicateBag(t *testing.T) {
	router := setupRouter()

	body := map[string]string{
		"qr_code":  "12345",
		"bag_type": "Parent",
		"status":   "Active",
	}
	jsonBody, _ := json.Marshal(body)

	// Register the first bag
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/create-bag", bytes.NewBuffer(jsonBody)))

	// Try registering the same bag again
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/create-bag", bytes.NewBuffer(jsonBody)))

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Bag with this QR code already exists")
}

func TestLinkBags(t *testing.T) {
	router := setupRouter()

	// Register Parent Bag
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/create-bag", bytes.NewBuffer([]byte(`{"qr_code":"PARENT123","bag_type":"Parent","status":"Active"}`))))

	// Register Child Bag
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/create-bag", bytes.NewBuffer([]byte(`{"qr_code":"CHILD123","bag_type":"Child","status":"Active"}`))))

	// Link Parent and Child
	linkBody := map[string]string{
		"parent_bag_qr_code": "PARENT123",
		"child_bag_qr_code":  "CHILD123",
	}
	jsonLinkBody, _ := json.Marshal(linkBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/link-bags", bytes.NewBuffer(jsonLinkBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Bags linked successfully")
}

func TestSearchBillByBag(t *testing.T) {
	router := setupRouter()

	// Register and link bag to a bill
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/create-bag", bytes.NewBuffer([]byte(`{"qr_code":"SEARCH123","bag_type":"Parent","status":"Active"}`))))
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/link-bag-to-bill", bytes.NewBuffer([]byte(`{"parent_bag_qr_code":"SEARCH123","bill_id":"BILL123"}`))))

	// Search for the bag
	req := httptest.NewRequest("GET", "/search-bill-by-bag?qr_code=SEARCH123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "BILL123")
}
