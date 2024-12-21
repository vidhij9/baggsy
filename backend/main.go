package main

import (
	"baggsy/controllers"
	"baggsy/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.InitDB()
	defer database.DB.Close()

	// Run migrations to ensure database tables are ready
	database.RunMigrations()

	// Set up the router
	r := gin.Default()

	// Bag-related endpoints
	r.POST("/create-bag", controllers.CreateBag)                    // Create a new bag
	r.GET("/bags", controllers.GetBagsWithPagination)               // Get paginated list of bags
	r.POST("/link-bags-to-sap-bill", controllers.LinkBagsToSAPBill) // Link bags to a bill

	// Bill-related endpoints
	r.POST("/create-bill", controllers.CreateBill) // Create a new bill and link parent bags
	r.GET("/bills", controllers.GetBills)          // Retrieve all bills

	// Start the server
	r.Run(":8080") // Listen and serve on port 8080
}
