package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/handlers"
	"baggsy/backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Rate limiting: 100 requests per second per IP
	limiter := rate.NewLimiter(rate.Every(time.Second/100), 100)
	r.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:80", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length", "X-Total-Count"}
	r.Use(cors.New(corsConfig))

	// Public routes
	r.POST("/login", handlers.LoginHandler)
	r.POST("/register", handlers.RegisterHandler)          // New account creation
	r.GET("/verify/:token", handlers.VerifyAccountHandler) // Email verification

	// Protected routes
	api := r.Group("/api").Use(middleware.AuthMiddleware())
	{
		api.POST("/register-parent", middleware.RestrictTo("employee", "admin"), handlers.RegisterParentHandler)
		api.POST("/register-child", middleware.RestrictTo("employee", "admin"), handlers.RegisterChildHandler)
		api.POST("/link-bags-to-bill", middleware.RestrictTo("employee", "admin"), handlers.LinkBagsToBillHandler)
		api.DELETE("/unlink-bag/:id", middleware.RestrictTo("employee", "admin"), handlers.UnlinkBagHandler)
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
