package routes

import (
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(OrderGroup *gin.RouterGroup) {
	OrderGroup.POST("/", services.CreateOrder)
	OrderGroup.GET("/", services.GetOrders)
	OrderGroup.GET("/:id", services.GetOrderByID)
	OrderGroup.PUT("/:id", services.UpdateOrder)
	OrderGroup.DELETE("/:id", services.DeleteOrder)

}
