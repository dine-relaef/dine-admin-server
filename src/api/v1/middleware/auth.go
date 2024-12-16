package middleware

import (
	"net/http"
	"strings"
	postgres "menu-server/src/config/database"
	"menu-server/src/models"
	"menu-server/src/utils"

	"github.com/gin-gonic/gin"
)

// Authenticate is a middleware function for verifying JWT tokens and user identity.
func Authenticate(c *gin.Context) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.Abort() // Stop further processing
		return
	}

	// Ensure the header contains "Bearer" token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	// Extract the token
	token := tokenParts[1]

	// Validate and extract user information from the token
	userID, role, err := utils.ValidateAndExtractToken(token, "ACCESS")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	// Verify the user exists in the database
	var user models.User
	if err := postgres.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	}

	// Attach user information to the request context
	c.Set("userID", userID)
	c.Set("role", role)

	// Continue processing the request
	c.Next()
}
