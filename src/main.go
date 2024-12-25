package main

import (
	docs "dine-server/docs"
	postgres "dine-server/src/config/database"
	"dine-server/src/config/env"
	"dine-server/src/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	r := gin.New()

	log.Println("Starting server...")
	log.Println("Environment: ", env.AppVar["ENVIRONMENT"])
	log.Println("App Name: ", env.PostgresDatabaseVar["DATABASE_URL"])
	if env.AppVar["ENVIRONMENT"] != "production" {
		docs.SwaggerInfo.Title = "Dine Server"
		docs.SwaggerInfo.Description = "Dine Server API documentation"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = env.AppVar["SERVER_HOST"]
		if env.AppVar["ENVIRONMENT"] == "testing" {
		docs.SwaggerInfo.Schemes = []string{"https"}
		}
		docs.SwaggerInfo.BasePath = "/"
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	r.Use(gin.Recovery())
	if env.AppVar["ENVIRONMENT"] == "development" {
		r.Use(gin.Logger())
	}

	postgres.InitDB()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},                          // Headers to expose
		AllowCredentials: true,                                                // Allow credentials like cookies
		MaxAge:           12 * 60 * 60,                                        // Cache preflight response for 12 hours
	}))

	setupRoutes(r)

	port := env.AppVar["PORT"]
	if port == "" {
		log.Println("PORT environment variable not set. Using default port 8080")
		port = "8080"
	}

	log.Printf("Starting %s on port %s in %s ", env.AppVar["APP_NAME"], port, env.AppVar["ENVIRONMENT"])
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

// setupHealthCheckRoute defines the health check endpoint
// @BasePath godoc
// @Summary  Health check endpoint
// @Schemes
// @Description Health check endpoint
// @Tags Health
// @Accept json
// @Produce json
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
	setupHealthCheckRoute(r)
	routes.V1Routes(r)
}
