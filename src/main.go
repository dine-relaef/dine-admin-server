package main

import (
	"log"
	docs "menu-server/docs"
	postgres "menu-server/src/config/database"
	"menu-server/src/config/env"
	"menu-server/src/routes"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	postgres.InitDB()

	r := gin.Default()
	if env.AppVar["ENVIRONMENT"] == "development" {
		docs.SwaggerInfo.Title = "Menu Server"
		docs.SwaggerInfo.Description = "Menu Server API documentation"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}

		docs.SwaggerInfo.BasePath = "/"
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

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
