// Main application entry point
package main

import (
	"baggsy/backend/controllers"
	"baggsy/backend/database"
	"log"
	"os"
	"strconv"
	"sync"

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

	// Set up the router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.Default())

	r.POST("/register-bag", controllers.RegisterBag)
	r.POST("/link-bags", controllers.LinkBags)
	r.POST("/link-bag-to-bill", controllers.LinkBagToBill)
	r.GET("/search-bill-by-bag", controllers.SearchBillByBag)
	r.GET("/linked-bags", controllers.GetLinkedBags)

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
