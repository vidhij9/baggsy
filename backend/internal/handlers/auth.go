package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("star-agri-seeds-secret")

func LoginHandler(c *gin.Context) {
	var creds struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		log.Printf("User not found: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if !user.Verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account not verified"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		log.Printf("Password mismatch: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiry
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func RegisterHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(req.Password) < 8 || !containsUpper(req.Password) || !containsNumber(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be 8+ chars with uppercase and number"})
		return
	}
	if req.Role != "employee" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'employee' or 'admin'"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	verificationToken := generateVerificationToken()
	user := models.User{
		Username:          req.Username,
		PasswordHash:      string(hash),
		Email:             req.Email,
		Role:              req.Role,
		VerificationToken: verificationToken,
		Verified:          false,
	}

	tx := db.DB.Begin()
	if !tx.Where("username = ?", req.Username).First(&models.User{}).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}
	if !tx.Where("email = ?", req.Email).First(&models.User{}).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	tx.Commit()

	// Simulate sending verification email (replace with actual email service)
	log.Printf("Verification link: https://baggsy-backend.up.railway.app/verify/%s\n", verificationToken)
	c.JSON(http.StatusOK, gin.H{"message": "Account created. Check email for verification link."})
}

func VerifyAccountHandler(c *gin.Context) {
	token := c.Param("token")
	var user models.User
	if err := db.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired verification token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify account"})
		}
		return
	}

	if user.Verified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account already verified"})
		return
	}

	if err := db.DB.Model(&user).Updates(map[string]interface{}{
		"verified":           true,
		"verification_token": "",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}

func containsUpper(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func generateVerificationToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
