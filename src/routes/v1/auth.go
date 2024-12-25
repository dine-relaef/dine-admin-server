package routes_v1

import (
	services "dine-server/src/api/v1/services/users"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(authGroup *gin.RouterGroup) {

	authGroup.POST("/register", services.RegisterUser)
	authGroup.POST("/login", services.LoginUser)
	authGroup.GET("/logout", services.LogoutUser)
	authGroup.GET("/refresh", services.RefreshToken)
	authGroup.GET("/google", services.GoogleLogin)
	authGroup.GET("/google/callback", services.GoogleCallback)

}
