package routes_v1

import (
	"menu-server/src/api/v1/middleware"
	services "menu-server/src/api/v1/services/subscriptions"

	"github.com/gin-gonic/gin"
)

func SetupSubscriptionRoutes(subscriptionRoutes *gin.RouterGroup) {

	// subscriptionRoutes.POST("/:payment_id", middleware.Authenticate, services.CreateSubscription)                                       // Create a subscription
	subscriptionRoutes.GET("/", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.GetAllSubscriptions)    // Get all subscriptions
	subscriptionRoutes.GET("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.GetSubscriptionByID) // Get a subscription by ID

}
