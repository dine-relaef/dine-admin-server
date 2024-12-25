package routes_v1

import (
	"dine-server/src/api/v1/middleware"
	"dine-server/src/api/v1/workflow"

	"github.com/gin-gonic/gin"
)

func SetupWorkflowRoutes(workflowGroup *gin.RouterGroup) {

	workflowGroup.POST("/plan/order-payment", middleware.Authenticate, middleware.RoleMiddleware([]string{"restaurant_admin"}), workflow.PlanOrderPayment)
	workflowGroup.GET("/plan/payment-subscription", middleware.Authenticate, middleware.RoleMiddleware([]string{"restaurant_admin"}), workflow.VerifyPaymentAndSubscription)

}
