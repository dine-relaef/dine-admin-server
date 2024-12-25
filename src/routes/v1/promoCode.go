package routes_v1

import (
	middleware "dine-server/src/api/v1/middleware"
	services_promocode "dine-server/src/api/v1/services/promocode"

	"github.com/gin-gonic/gin"
)

func SetupPromoCodeRoutes(promoCodeGroup *gin.RouterGroup) {
	promoCodeGroup.POST("/dine", middleware.Authenticate, middleware.RoleMiddleware([]string{"admin"}), services_promocode.CreateDinePromoCode)
}

func DinePromoCodeRoutes(promoCodeDineGroup *gin.RouterGroup) {

	promoCodeDineGroup.POST("/", services_promocode.CreateDinePromoCode)

}
