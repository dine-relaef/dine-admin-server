package routes_v1

import (
	"menu-server/src/api/v1/middleware"
	services_orders "menu-server/src/api/v1/services/orders"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(orderGroup *gin.RouterGroup) {

	orderGroup.POST("", services_orders.CreateOrder)
	orderGroup.GET("", services_orders.ListOrders)
	orderGroup.GET("/:id", services_orders.GetOrder)
	orderGroup.PUT("/:id/status", services_orders.UpdateOrderStatus)
	orderGroup.POST("/:id/cancel", services_orders.CancelOrder)

	orderGroup.POST("/dine", middleware.Authenticate, middleware.RoleMiddleware([]string{"restaurant_admin"}), services_orders.CreateDineOrder)
	orderGroup.GET("/dine/all", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_orders.GetDineOrders)
	orderGroup.GET("/dine/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_orders.GetDineOrderByID)
	orderGroup.GET("/dine", middleware.Authenticate, middleware.RoleMiddleware([]string{"restaurant_admin"}), services_orders.GetDineOrderByUsers)

}
