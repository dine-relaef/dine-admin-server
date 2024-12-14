package routes

import (
	controllers "menu-server/src/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.Engine) {
	PaymentGroup := router.Group("/payments")
	{
		PaymentGroup.POST("/", controllers.CreatePayment)
		PaymentGroup.GET("/", controllers.GetPayments)
		PaymentGroup.GET("/:id", controllers.GetPaymentByID)
		PaymentGroup.PUT("/:id", controllers.UpdatePayment)
		// PaymentGroup.DELETE("/:id", controllers.DeletePayment)
	}
}