package routes

import (
	controllers "menu-server/src/controllers/v1"
	"github.com/gin-gonic/gin"
)

func SetupSubscriptionRoutes(r *gin.Engine) {
	subscriptionRoutes := r.Group("/subscriptions")
	{
		subscriptionRoutes.POST("/", controllers.CreateSubscription)          // Create a subscription
		subscriptionRoutes.GET("/", controllers.GetAllSubscriptions)          // Get all subscriptions
		subscriptionRoutes.GET("/:id", controllers.GetSubscriptionByID)      // Get a subscription by ID
		subscriptionRoutes.PUT("/:id", controllers.UpdateSubscription)       // Update a subscription by ID
		subscriptionRoutes.DELETE("/:id", controllers.DeleteSubscription)    // Delete a subscription by ID
	}
}
