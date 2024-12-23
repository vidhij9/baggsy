package tests

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"../backend/controllers"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCreateBagAPI(t *testing.T) {
// 	router := gin.Default()
// 	router.POST("/create-bag", controllers.CreateBag)

// 	t.Run("Valid Bag Creation", func(t *testing.T) {
// 		body := `{"id":"bag1", "qr_code":"testQR1", "bag_type":"Parent", "status":"Active"}`
// 		req, _ := http.NewRequest("POST", "/create-bag", bytes.NewBuffer([]byte(body)))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusCreated, resp.Code)
// 	})

// 	t.Run("Duplicate Bag Creation", func(t *testing.T) {
// 		body := `{"id":"bag1", "qr_code":"testQR1", "bag_type":"Parent", "status":"Active"}`
// 		req, _ := http.NewRequest("POST", "/create-bag", bytes.NewBuffer([]byte(body)))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusConflict, resp.Code)
// 	})
// }
