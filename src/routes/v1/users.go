package routes

import (
	"menu-server/src/api/v1/middleware"
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up the routes for the user resource
// @Summary Set up user routes
func SetupUserRoutes(userGroup *gin.RouterGroup) {

	
	userGroup.GET("/", middleware.RoleMiddleware("admin"), services.GetUsers)
	userGroup.GET("/:id", services.GetUserByID)
	userGroup.PUT("/:id", services.UpdateUser)
	userGroup.DELETE("/:id", services.DeleteUser)

}
