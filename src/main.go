package main

import (
	"log"
	docs "menu-server/docs"
	postgres "menu-server/src/config/database"
	"menu-server/src/config/env"
	routes "menu-server/src/routes/v1"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize the database
	postgres.InitDB()

	// Create a new Gin router
	r := gin.Default()
	if env.AppVar["ENVIRONMENT"] == "development" {
		// Swagger documentation
		docs.SwaggerInfo.Title = "Menu Server"
		docs.SwaggerInfo.Description = "Menu Server API documentation"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}

		docs.SwaggerInfo.BasePath = "/"
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

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

// Response represents the successful API response
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// setupHealthCheckRoute defines the health check endpoint
// @BasePath godoc
// @Summary  Health check endpoint
// @Schemes
// @Description Health check endpoint
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router / [get]
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

// V1Routes sets up all version 1 routes
func V1Routes(r *gin.Engine) {
	v1 := r.Group("/api/v1") // Create a /api/v1 route group

	// @Tags Users
	routes.SetupAuthRoutes(v1.Group("/auth"))
	routes.SetupUserRoutes(v1.Group("/users"))
	// @Tags Subscriptions
	routes.SetupSubscriptionRoutes(v1.Group("/subscriptions"))
	routes.SetupRestaurantRoutes(v1.Group("/restaurants"))
	routes.SetupPlanRoutes(v1.Group("/plans"))
	routes.SetupPaymentRoutes(v1.Group("/payments"))
	routes.SetupOrderRoutes(v1.Group("/orders"))
	routes.SetupMenuRoutes(v1.Group("/menus"))
}
