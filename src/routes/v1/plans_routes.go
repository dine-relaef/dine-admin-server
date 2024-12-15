package routes

import (
	controllers "menu-server/src/controllers/v1"
	"github.com/gin-gonic/gin"
	middleware "menu-server/src/api/v1"
)

func SetupPlanRoutes(router *gin.Engine) {
	PlanGroup := router.Group("/plans")
	{
		PlanGroup.POST("/", middleware.RoleMiddleware("admin"), controllers.CreatePlan)
		PlanGroup.GET("/", controllers.GetPlans)
		PlanGroup.GET("/:id", controllers.GetPlanByID)
		PlanGroup.PUT("/:id", middleware.RoleMiddleware("admin"), controllers.UpdatePlan)
		PlanGroup.DELETE("/:id", middleware.RoleMiddleware("admin"), controllers.DeletePlan)
	}
}