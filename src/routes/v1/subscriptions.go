package routes_v1

import (
	"menu-server/src/api/v1/middleware"
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupSubscriptionRoutes(subscriptionRoutes *gin.RouterGroup) {

	subscriptionRoutes.POST("/:payment_id", middleware.Authenticate, services.CreateSubscription) // Create a subscription
	subscriptionRoutes.GET("/", services.GetAllSubscriptions)                                          // Get all subscriptions
	subscriptionRoutes.GET("/:id", services.GetSubscriptionByID)                                       // Get a subscription by ID
	subscriptionRoutes.PUT("/:id", services.UpdateSubscription)                                        // Update a subscription by ID
	subscriptionRoutes.DELETE("/:id", services.DeleteSubscription)                                     // Delete a subscription by ID

}
