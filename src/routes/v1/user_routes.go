package routes

import (
	controllers "menu-server/src/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/", controllers.GetUsers)
		userGroup.GET("/:id", controllers.GetUserByID)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}
}