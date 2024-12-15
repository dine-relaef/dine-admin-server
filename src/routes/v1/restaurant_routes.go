package routes

import (
	controllers "menu-server/src/controllers/v1"
	"github.com/gin-gonic/gin"
	middleware "menu-server/src/api/v1"
)

func SetupRestaurantRoutes(r *gin.Engine) {
	RestaurantRoutes := r.Group("/restaurant")
	{
		RestaurantRoutes.POST("/", middleware.RoleMiddleware("admin"), controllers.CreateRestaurant)          // Create a Restaurant
		RestaurantRoutes.GET("/", controllers.GetAllRestaurants)          // Get all Restaurants
		RestaurantRoutes.GET("/:id", controllers.GetRestaurantByID)      // Get a Restaurant by ID
		RestaurantRoutes.PUT("/:id", controllers.UpdateRestaurant)       // Update a Restaurant by ID
		RestaurantRoutes.DELETE("/:id", middleware.RoleMiddleware("admin"), controllers.DeleteRestaurant)    // Delete a Restaurant by ID
	}
}
