package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func contains(slice []string, item interface{}) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RoleMiddleware(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user_id from query parameters
		role, _ := c.Get("role")
		if !contains(requiredRoles, role) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
