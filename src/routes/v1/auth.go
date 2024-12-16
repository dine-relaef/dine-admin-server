package routes

import (
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(authGroup *gin.RouterGroup) {

	authGroup.POST("/register", services.RegisterUser)
	authGroup.POST("/login", services.LoginUser)

}
