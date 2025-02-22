// package main

// import (
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// 	"golang.org/x/crypto/bcrypt"
// )

// var db *gorm.DB
// var jwtSecret = []byte("secret-key") // Replace with env var in production

// func main() {
// 	// Connect to PostgreSQL
// 	var err error
// 	db, err = gorm.Open("postgres", "host=localhost user=baggsy password=baggsy dbname=baggsy sslmode=disable")
// 	if err != nil {
// 		panic("Failed to connect to database")
// 	}
// 	defer db.Close()
// 	db.AutoMigrate(&User{}, &Bag{}, &Link{})

// 	// Initialize Gin router
// 	r := gin.Default()

// 	// Public routes
// 	r.POST("/login", loginHandler)

// 	// Protected routes
// 	auth := r.Group("/api").Use(authMiddleware())
// 	{
// 		auth.POST("/register-bag", registerBagHandler)
// 		auth.POST("/link-bag-to-bill", linkBagToBillHandler)
// 		auth.GET("/search/bag", searchBagHandler)
// 		auth.GET("/search/bill", searchBillHandler)
// 	}

// 	r.Run(":8080")
// }

// // Models
// type User struct {
// 	ID           uint   `gorm:"primary_key"`
// 	Username     string `gorm:"unique"`
// 	PasswordHash string
// 	Role         string
// }

// type Bag struct {
// 	ID         uint   `gorm:"primary_key"`
// 	QRCode     string `gorm:"unique"`
// 	Type       string
// 	ChildCount int
// 	ParentID   *uint
// 	Linked     bool
// 	CreatedAt  time.Time
// }

// type Link struct {
// 	ID        uint `gorm:"primary_key"`
// 	ParentID  uint `gorm:"unique"`
// 	BillID    string
// 	CreatedAt time.Time
// }

// // Middleware
// func authMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")[7:] // Remove "Bearer "
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return jwtSecret, nil
// 		})
// 		if err != nil || !token.Valid {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
// 			return
// 		}
// 		c.Next()
// 	}
// }

// // Handlers
// func loginHandler(c *gin.Context) {
// 	var creds struct{ Username, Password string }
// 	c.BindJSON(&creds)
// 	var user User
// 	if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
// 		c.JSON(401, gin.H{"error": "Invalid credentials"})
// 		return
// 	}
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
// 		c.JSON(401, gin.H{"error": "Invalid credentials"})
// 		return
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userId": user.ID,
// 		"role":   user.Role,
// 		"exp":    time.Now().Add(time.Hour * 1).Unix(),
// 	})
// 	tokenString, _ := token.SignedString(jwtSecret)
// 	c.JSON(200, gin.H{"token": tokenString})
// }

// func registerBagHandler(c *gin.Context) {
// 	var bag Bag
// 	c.BindJSON(&bag)
// 	if db.Where("qr_code = ?", bag.QRCode).First(&Bag{}).RecordNotFound() {
// 		if bag.Type == "child" && bag.ParentID != nil {
// 			var parent Bag
// 			db.First(&parent, *bag.ParentID)
// 			// Convert RowsAffected to int for comparison
// 			currentChildCount := int(db.Where("parent_id = ?", *bag.ParentID).Find(&[]Bag{}).RowsAffected)
// 			if parent.ChildCount <= currentChildCount {
// 				c.JSON(400, gin.H{"error": "Parent capacity reached"})
// 				return
// 			}
// 		}
// 		db.Create(&bag)
// 		c.JSON(200, bag)
// 	} else {
// 		c.JSON(400, gin.H{"error": "QR code already registered"})
// 	}
// }

// func linkBagToBillHandler(c *gin.Context) {
// 	var link Link
// 	c.BindJSON(&link)
// 	var bag Bag
// 	if db.First(&bag, link.ParentID).Error != nil || bag.Linked {
// 		c.JSON(400, gin.H{"error": "Bag not found or already linked"})
// 		return
// 	}
// 	db.Create(&link)
// 	db.Model(&bag).Update("linked", true)
// 	c.JSON(200, link)
// }

// func searchBagHandler(c *gin.Context) {
// 	qr := c.Query("qr")
// 	var bag Bag
// 	if db.Where("qr_code = ?", qr).First(&bag).Error != nil {
// 		c.JSON(404, gin.H{"error": "Bag not found"})
// 		return
// 	}
// 	if bag.Type == "child" && bag.ParentID != nil {
// 		var parent Bag
// 		var link Link
// 		db.First(&parent, *bag.ParentID)
// 		db.Where("parent_id = ?", *bag.ParentID).First(&link)
// 		c.JSON(200, gin.H{"bag": bag, "parent": parent, "bill": link.BillID})
// 	} else {
// 		var link Link
// 		db.Where("parent_id = ?", bag.ID).First(&link)
// 		c.JSON(200, gin.H{"bag": bag, "bill": link.BillID})
// 	}
// }

// func searchBillHandler(c *gin.Context) {
// 	billID := c.Query("bill")
// 	var links []Link
// 	db.Where("bill_id = ?", billID).Find(&links)
// 	var bags []Bag
// 	for _, link := range links {
// 		var bag Bag
// 		db.First(&bag, link.ParentID)
// 		bags = append(bags, bag)
// 	}
// 	c.JSON(200, bags)
// }

package main

import (
	"fmt"
	"log"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/handlers"
	"baggsy/backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Setup Gin router
	r := gin.New()      // Use New() instead of Default() to add middleware explicitly
	r.Use(gin.Logger()) // Add logging explicitly
	r.Use(gin.Recovery())

	// CORS configuration - Apply before routes
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:80", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(corsConfig))

	// Public routes
	r.POST("/login", handlers.LoginHandler)

	// Protected routes with auth middleware
	api := r.Group("/api").Use(middleware.AuthMiddleware())
	{
		// Employee-accessible routes
		api.POST("/register-parent", middleware.RestrictTo("employee", "admin"), handlers.RegisterParentHandler)
		api.POST("/register-child", middleware.RestrictTo("employee", "admin"), handlers.RegisterChildHandler)
		api.POST("/link-bags-to-bill", middleware.RestrictTo("employee", "admin"), handlers.LinkBagsToBillHandler)
		api.DELETE("/unlink-bag/:id", middleware.RestrictTo("employee", "admin"), handlers.UnlinkBagHandler)

		// Admin-only routes
		api.GET("/bags", middleware.RestrictTo("admin"), handlers.ListBagsHandler)
		api.GET("/unlinked-parents", middleware.RestrictTo("admin"), handlers.ListUnlinkedParentsHandler)
		api.GET("/child-bags/:parentQR", middleware.RestrictTo("admin"), handlers.FindChildBagsByParentQRHandler)
		api.GET("/bills", middleware.RestrictTo("admin"), handlers.ListBillsHandler)
		api.GET("/bill/:billID", middleware.RestrictTo("admin"), handlers.SearchBillByNumberHandler)
		api.GET("/bag/:qr", middleware.RestrictTo("admin"), handlers.SearchBagByQRHandler)
	}

	fmt.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
