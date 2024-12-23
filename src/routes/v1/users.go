package routes_v1

import (
	"menu-server/src/api/v1/middleware"
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up the routes for the user resource
// @Summary Set up user routes
func SetupUserRoutes(userGroup *gin.RouterGroup) {

	userGroup.GET("/get-all", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.GetAllUsers)
	userGroup.GET("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.GetUserByID)
	userGroup.PUT("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.UpdateUser)
	userGroup.DELETE("/:id", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services.DeleteUser)

	userGroup.GET("/", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin", "restaurant_admin"}), services.GetUser)
	userGroup.PUT("/", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin", "restaurant_admin"}), services.UpdateUserByUser)

}
