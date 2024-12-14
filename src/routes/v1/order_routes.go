package routes

import (
	controllers "menu-server/src/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine) {
	OrderGroup := router.Group("/Orders")
	{
		OrderGroup.POST("/", controllers.CreateOrder)
		OrderGroup.GET("/", controllers.GetOrders)
		OrderGroup.GET("/:id", controllers.GetOrderByID)
		OrderGroup.PUT("/:id", controllers.UpdateOrder)
		OrderGroup.DELETE("/:id", controllers.DeleteOrder)
	}
}
