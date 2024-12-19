package routes

import (
	middleware "menu-server/src/api/v1/middleware"
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupRestaurantRoutes(RestaurantRoutes *gin.RouterGroup) {
	RestaurantRoutes.GET("/get-all", services.GetAllRestaurants)
	RestaurantRoutes.GET("/", middleware.Authenticate, services.GetRestaurants)                                                                         // Get all Restaurants
	RestaurantRoutes.GET("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin", "restaurant_admin"}), services.GetRestaurantByID) // Get a Restaurant by ID

	RestaurantRoutes.POST("/", middleware.Authenticate, services.CreateRestaurant) // Create a Restaurant

	RestaurantRoutes.PUT("/:id", middleware.Authenticate, services.UpdateRestaurant) // Update a Restaurant by ID

	RestaurantRoutes.DELETE("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin", "restaurant_admin"}), services.DeleteRestaurant) // Delete a Restaurant by ID

}
