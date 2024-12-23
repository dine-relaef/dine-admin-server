package routes_v1

import (
	"menu-server/src/api/v1/middleware"
	services_payments "menu-server/src/api/v1/services/payments"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(PaymentGroup *gin.RouterGroup) {

	PaymentGroup.POST("/dine/:order_id", middleware.Authenticate, services_payments.CreateDinePayment)
	PaymentGroup.GET("/dine/callback", services_payments.PaymentCallback)
	// PaymentGroup.GET("/", services_payments.GetPayments)
	// PaymentGroup.GET("/:id", services_payments.GetPaymentByID)

}
