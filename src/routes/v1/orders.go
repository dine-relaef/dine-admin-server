package routes

import (
	"menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(orderGroup *gin.RouterGroup) {

	orderGroup.POST("", services.CreateOrder)
	orderGroup.GET("", services.ListOrders)
	orderGroup.GET("/:id", services.GetOrder)
	orderGroup.PUT("/:id/status", services.UpdateOrderStatus)
	orderGroup.POST("/:id/cancel", services.CancelOrder)

}
