package routes

import (
	middleware "menu-server/src/api/v1/middleware"
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupRestaurantRoutes(RestaurantRoutes *gin.RouterGroup) {

	RestaurantRoutes.POST("/", middleware.RoleMiddleware("admin"), services.CreateRestaurant)      // Create a Restaurant
	RestaurantRoutes.GET("/", services.GetAllRestaurants)                                          // Get all Restaurants
	RestaurantRoutes.GET("/:id", services.GetRestaurantByID)                                       // Get a Restaurant by ID
	RestaurantRoutes.PUT("/:id", services.UpdateRestaurant)                                        // Update a Restaurant by ID
	RestaurantRoutes.DELETE("/:id", middleware.RoleMiddleware("admin"), services.DeleteRestaurant) // Delete a Restaurant by ID

}
