package main

import (
	"log"
	"os"

	"database"
	"models"

	"./controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	db := database.ConnectDB()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get DB instance: %v", err)
		}
		sqlDB.Close()
	}()

	// Run migrations
	database.RunMigrations(db)
	models.RegisterModels(db)

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	router.POST("/create-bag", controllers.CreateBag)
	router.POST("/link-bags-to-sap-bill", controllers.LinkBagsToSAPBill)
	router.GET("/bags", controllers.GetBagsWithPagination)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
