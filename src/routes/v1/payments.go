package routes_v1

import (
	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(PaymentGroup *gin.RouterGroup) {

	dinePaymentRoutes(PaymentGroup.Group("/dine"))
	// PaymentGroup.GET("/", services_payments.GetPayments)
	// PaymentGroup.GET("/:id", services_payments.GetPaymentByID)

}

func dinePaymentRoutes(PaymentDineGroup *gin.RouterGroup) {

	// PaymentDineGroup.GET("/callback", services_payments.PaymentCallback)
	// PaymentGroup.GET("/", services_payments.GetPayments)
	// PaymentGroup.GET("/:id", services_payments.GetPaymentByID)

}
