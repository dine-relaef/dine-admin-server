package middleware

import (
	"net/http"
	models "menu-server/src/models"
	postgres "menu-server/src/config/database"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user_id from query parameters
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
			c.Abort()
			return
		}

		var user models.User
		if err := postgres.DB.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Check if the user's role matches the required role
		if user.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied, insufficient role"})
			c.Abort()
			return
		}

		// Add user_id and role to context for further use
		c.Set("user_id", userID)
		c.Set("role", user.Role)
		c.Next()
	}
}