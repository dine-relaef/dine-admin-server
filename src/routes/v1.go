package routes

import (
	routes_v1 "dine-server/src/routes/v1"

	"github.com/gin-gonic/gin"
)

// V1Routes sets up all version 1 routes
func V1Routes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	routes_v1.SetupAuthRoutes(v1.Group("/auth"))
	routes_v1.SetupUserRoutes(v1.Group("/users"))
	routes_v1.SetupSubscriptionRoutes(v1.Group("/subscriptions"))
	routes_v1.SetupRestaurantRoutes(v1.Group("/restaurants"))
	routes_v1.SetupPlanRoutes(v1.Group("/plans"))
	routes_v1.SetupPaymentRoutes(v1.Group("/payments"))
	routes_v1.SetupOrderRoutes(v1.Group("/orders"))
	routes_v1.SetupMenuRoutes(v1.Group("/:restaurant_id/menus"))
	routes_v1.SetupPromoCodeRoutes(v1.Group("/promo-code"))
	routes_v1.SetupWorkflowRoutes(v1.Group("/workflow"))
}
