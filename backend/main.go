package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"baggsy/backend/controllers"
	"baggsy/backend/database"

	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var mu sync.Mutex // Global mutex for thread-safe operations

func main() {
	// Initialize the database
	database.Connect()

	defer func() {
		mu.Lock()
		defer mu.Unlock()
		if err := database.Close(); err != nil {
			log.Fatalf("Error closing the database connection: %v", err)
		}
	}()

	// Run migrations
	mu.Lock()
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Migration execution error: %v", err)
	}
	mu.Unlock()

	// Set up the GIN router
	r := gin.Default()

	// Create a ginprom Prometheus middleware
	gp := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("baggsy"),
		ginprom.Path("/metrics"),
	)

	// Attach the middleware
	r.Use(gp.Instrument())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Add CORS middleware
	r.Use(cors.Default())

	// Optional: Rate-limit middleware (example with a simple token bucket).
	r.Use(rateLimitMiddleware(100)) // e.g. 100 requests/second per instance

	r.POST("/register-bag", controllers.RegisterBag)
	r.POST("/link-child-bag", controllers.LinkChildBag)
	r.POST("/link-bag-to-bill", controllers.LinkBagToBill)
	r.GET("/search-bill", controllers.SearchBag)

	// Use a configurable port from environment variables
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Println("SERVER_PORT is not set, using default port 8080")
		port = "8080"
	}

	if !isValidPort(port) {
		log.Fatalf("Invalid port specified: %s", port)
	}

	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// isValidPort validates if the given port is within the valid range
func isValidPort(port string) bool {
	const minPort, maxPort = 1, 65535
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return portInt >= minPort && portInt <= maxPort
}

// Example rate-limit middleware (very simplistic):
func rateLimitMiddleware(rps int) gin.HandlerFunc {
	// In production, you might use a more robust library or Redis-based approach
	bucket := make(chan struct{}, rps)

	// Fill the bucket
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			for i := 0; i < rps-len(bucket); i++ {
				bucket <- struct{}{}
			}
		}
	}()

	return func(c *gin.Context) {
		select {
		case <-bucket:
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}
	}
}
