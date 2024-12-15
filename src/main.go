package main

import (
	"log"
	"os"
	postgres "menu-server/src/config/database"
	routes "menu-server/src/routes/v1"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	postgres.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Health check route
	setupHealthCheckRoute(r)

	// Initialize application routes
	setupRoutes(r)

	// Determine the port from environment variables (default: 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupHealthCheckRoute defines the health check endpoint
func setupHealthCheckRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Jai Mata Di",
			"status":  "ok",
		})
	})
}

// setupRoutes initializes all application routes
func setupRoutes(r *gin.Engine) {
	routes.SetupSubscriptionRoutes(r)
	routes.SetupUserRoutes(r)
	routes.SetupRestaurantRoutes(r)
	routes.SetupPlanRoutes(r)
	routes.SetupPaymentRoutes(r)
	routes.SetupOrderRoutes(r)
	routes.SetupMenuRoutes(r)
}
