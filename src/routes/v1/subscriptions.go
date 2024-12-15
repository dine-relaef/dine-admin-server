package routes

import (
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupSubscriptionRoutes(subscriptionRoutes *gin.RouterGroup) {

	subscriptionRoutes.POST("/", services.CreateSubscription)      // Create a subscription
	subscriptionRoutes.GET("/", services.GetAllSubscriptions)      // Get all subscriptions
	subscriptionRoutes.GET("/:id", services.GetSubscriptionByID)   // Get a subscription by ID
	subscriptionRoutes.PUT("/:id", services.UpdateSubscription)    // Update a subscription by ID
	subscriptionRoutes.DELETE("/:id", services.DeleteSubscription) // Delete a subscription by ID

}
