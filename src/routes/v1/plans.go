package routes_v1

import (
	middleware "dine-server/src/api/v1/middleware"
	services_plan "dine-server/src/api/v1/services/plans"

	"github.com/gin-gonic/gin"
)

func SetupPlanRoutes(PlanGroup *gin.RouterGroup) {

	//Open routes
	PlanGroup.GET("/", services_plan.GetPlans)
	PlanGroup.GET("/:id", services_plan.GetPlanByID)

	//Only accessible by admin
	PlanGroup.POST("/", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.CreatePlan)
	PlanGroup.GET("/all", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.GetAllPlans)
	PlanGroup.PUT("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.UpdatePlan)
	PlanGroup.DELETE("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.DeletePlan)
	PlanGroup.PUT("/add-feature", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.AddPlanFeature)
	PlanGroup.PUT("/remove-feature", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.RemovePlanFeature)
	// Plan Features
	PlanGroup.POST("/feature", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.CreatePlanFeatures)
	PlanGroup.GET("/feature", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.GetAllPlanFeatures)
	PlanGroup.PUT("/feature/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.UpdatePlanFeatures)
	PlanGroup.DELETE("/feature/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_plan.DeletePlanFeatures)

}
