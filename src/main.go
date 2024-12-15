package main

import (
	"log"
	postgres "menu-server/src/config/database"
	"menu-server/src/config/env"
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
	setupRoutes(r)
	// Determine the port from environment variables (default: 8080)
	port := env.AppVar["PORT"]
	if port == "" {
		log.Println("PORT environment variable not set. Using default port 8080")
		port = "8080"
	}

	// Start the server
	log.Printf("Starting %s on port %s in %s ", env.AppVar["APP_NAME"], port, env.AppVar["ENVIRONMENT"])
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// Initialize application routes

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

	V1Routes(r)
}

func V1Routes(r *gin.Engine) {
	v1 := r.Group("/api/v1") // Create a /api/v1 route group

	routes.SetupUserRoutes(v1.Group("/users"))
	routes.SetupSubscriptionRoutes(v1.Group("/subscriptions"))
	routes.SetupRestaurantRoutes(v1.Group("/restaurants"))
	routes.SetupPlanRoutes(v1.Group("/plans"))
	routes.SetupPaymentRoutes(v1.Group("/payments"))
	routes.SetupOrderRoutes(v1.Group("/orders"))
	routes.SetupMenuRoutes(v1.Group("/menus"))
}

