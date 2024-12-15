package routes

import (
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(userGroup *gin.RouterGroup) {

	userGroup.POST("/", services.CreateUser)
	userGroup.GET("/", services.GetUsers)
	userGroup.GET("/:id", services.GetUserByID)
	userGroup.PUT("/:id", services.UpdateUser)
	userGroup.DELETE("/:id", services.DeleteUser)

}
