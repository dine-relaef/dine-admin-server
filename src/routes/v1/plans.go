package routes

import (
	middleware "menu-server/src/api/v1/middleware"
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupPlanRoutes(PlanGroup *gin.RouterGroup) {

	PlanGroup.POST("/", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.CreatePlan)
	PlanGroup.GET("/", services.GetPlans)
	PlanGroup.GET("/:id", services.GetPlanByID)
	PlanGroup.PUT("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.UpdatePlan)
	PlanGroup.DELETE("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.DeletePlan)

}
