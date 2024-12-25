package routes_v1

import (
	"dine-server/src/api/v1/middleware"
	services_orders "dine-server/src/api/v1/services/orders"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(orderGroup *gin.RouterGroup) {

	restaurantOrderRoutes(orderGroup.Group("/restaurant"))
	dineOrderRoutes(orderGroup.Group("/dine"))

}

func restaurantOrderRoutes(orderRestaurantGroup *gin.RouterGroup) {

	orderRestaurantGroup.POST("/", services_orders.CreateOrder)
	orderRestaurantGroup.GET("/", services_orders.ListOrders)
	orderRestaurantGroup.GET("/:id", services_orders.GetOrder)
	orderRestaurantGroup.PUT("/:id/status", services_orders.UpdateOrderStatus)
	orderRestaurantGroup.POST("/:id/cancel", services_orders.CancelOrder)

}

func dineOrderRoutes(orderDineGroup *gin.RouterGroup) {

	orderDineGroup.GET("/all", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_orders.GetDineOrders)
	orderDineGroup.GET("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_orders.GetDineOrderByID)
	orderDineGroup.GET("/", middleware.Authenticate, middleware.RoleMiddleware([]string{"restaurant_admin"}), services_orders.GetDineOrderByUsers)

}
