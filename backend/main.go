// Main application entry point
package main

import (
	"baggsy/backend/controllers"
	"baggsy/backend/database"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.InitDBWithRetry()
	defer func() {
		if err := database.DB.Close(); err != nil {
			log.Fatalf("Error closing the database connection: %v", err)
		}
	}()

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Migration execution error: %v", err)
	}

	// Set up the router
	r := gin.Default()

	r.POST("/create-bag", controllers.CreateBag)
	r.POST("/link-bags-to-sap-bill", controllers.LinkBagsToSAPBill)
	r.GET("/bags", controllers.GetBagsWithPagination)

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
