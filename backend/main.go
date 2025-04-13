package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"baggsy/backend/internal/db"
	"baggsy/backend/internal/handlers"
	"baggsy/backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize Database connection (GORM)
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB := database.DB()
	defer sqlDB.Close()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Global Rate Limiter: 100 req/sec
	limiter := rate.NewLimiter(rate.Every(time.Second/100), 100)
	r.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	})

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://baggsy-frontend.up.railway.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// CORS Preflight globally handled
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	// Public endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	r.POST("/login", handlers.LoginHandler)
	r.POST("/register", handlers.RegisterHandler)
	r.GET("/verify/:token", handlers.VerifyAccountHandler)

	// Protected API routes
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

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
