package routes

import (
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(PaymentGroup *gin.RouterGroup) {

	PaymentGroup.POST("/", services.CreatePayment)
	PaymentGroup.GET("/", services.GetPayments)
	PaymentGroup.GET("/:id", services.GetPaymentByID)
	PaymentGroup.PUT("/:id", services.UpdatePayment)
	// PaymentGroup.DELETE("/:id", services.DeletePayment)

}
