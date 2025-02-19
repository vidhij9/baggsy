package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"

	"baggsy/backend/controllers"
	"baggsy/backend/database"
)

// Global database handle (opened in main and closed on exit)
var db *sql.DB
var limiters = make(map[string]*rate.Limiter) // for rate limiting by client IP

// RateLimiter struct to track requests per client
type RateLimiter struct {
	mu        sync.RWMutex
	requests  map[string]int       // request counts per client (IP)
	lastReset map[string]time.Time // last reset time per client
	limit     int                  // max requests per window
	window    time.Duration        // time window for rate limit
}

var rateLimiter = RateLimiter{
	requests:  make(map[string]int),
	lastReset: make(map[string]time.Time),
	limit:     100,             // e.g., 100 requests
	window:    1 * time.Minute, // per 1 minute window
}

// Handler for the /search-bill route
func searchBillHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	log.Printf("Searching for bill: %q", query)
	// ... perform search logic (using `db` as needed) ...
	w.Write([]byte("Results for query: " + query))
}

// rateLimitMiddleware limits each client to a fixed request rate to prevent abuse.
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		// Initialize a limiter for new IPs
		if _, exists := limiters[ip]; !exists {
			limiters[ip] = rate.NewLimiter(1, 5) // 1 req/sec, burst capacity 5
		}
		limiter := limiters[ip]

		// Allow or reject the request based on the rate limit
		if !limiter.Allow() {
			// Rate limit exceeded
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}

		c.Next() // within rate limit, proceed to handler
	}
}

// loggingMiddleware logs the incoming HTTP requests and response details.
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next() // process request

		// After the request is handled, record duration and other details
		duration := time.Since(startTime)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()

		log.Printf("Request: %s %s from %s | Status: %d | Duration: %v",
			method, path, clientIP, status, duration) //  [oai_citation_attribution:5‡codesignal.com](https://codesignal.com/learn/courses/enhancing-our-todo-app/lessons/implementing-request-logging-middleware#:~:text=17%20%20%20%20,method%2C%20path%2C%20clientIP%2C%20statusCode%2C%20duration)
	}
}

func main() {
	// Initialize database and run migrations
	db, err := database.ConnectDB("dsn") // (Assuming ConnectDB opens a GORM *gorm.DB)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	database.RunMigrations(db) // run auto-migrations to create/update tables

	// Set up Gin router and middleware
	r := gin.New()
	r.Use(loggingMiddleware())   // custom request logging
	r.Use(rateLimitMiddleware()) // rate limiting to prevent abuse
	r.Use(gin.Recovery())        // panic recovery middleware

	// Bag Management Endpoints
	r.POST("/register-bag", controllers.RegisterBag)   // Register a new bag (Parent or Child)
	r.GET("/bags", controllers.ListBags)               // List all bags (with pagination & filters)
	r.GET("/bags/:qrCode", controllers.GetBagByQRCode) // Get details of a specific bag by QR code

	// Child-Parent Linking Endpoints
	r.POST("/link-child-bag", controllers.LinkChildBag)           // Link a child bag to a parent bag
	r.POST("/unlink-child-bag", controllers.UnlinkChildBag)       // Unlink a child bag from its parent
	r.GET("/parent/:parentId/children", controllers.GetChildBags) // List all child bags under a given parent

	// Parent-Bill Linking Endpoints
	r.POST("/link-bag-to-bill", controllers.LinkBagToBill)         // Link a parent bag to a bill
	r.POST("/unlink-bag-from-bill", controllers.UnlinkBagFromBill) // Unlink a parent bag from a bill
	r.GET("/bill/:billId/parents", controllers.GetParentBags)      // List all parent bags linked to a given bill

	// Search & Lookup Endpoints
	r.GET("/search-bag", controllers.SearchBag)   // Search for a bag by QR code (query param q)
	r.GET("/search-bill", controllers.SearchBill) // Search for a bill by ID (query param q)

	// Health & Monitoring Endpoints
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// Expose Prometheus metrics at /metrics using promhttp handler
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) //  [oai_citation_attribution:2‡stackoverflow.com](https://stackoverflow.com/questions/65608610/how-to-use-gin-as-a-server-to-write-prometheus-exporter-metrics#:~:text=Use%20gin%20wrapper)

	// Start the HTTP server
	if err := r.Run(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
