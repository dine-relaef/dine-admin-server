package middleware

import (
	postgres "dine-server/src/config/database"
	models_user "dine-server/src/models/users"
	"dine-server/src/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authenticate is a middleware function for verifying JWT tokens and user identity.
func Authenticate(c *gin.Context) {
	// Get token age from environment variables
	accessTokenAge := utils.ParseDuration(os.Getenv("ACCESS_TOKEN_AGE"), 3600)    // Default 3600 seconds (1 hour)
	refreshTokenAge := utils.ParseDuration(os.Getenv("REFRESH_TOKEN_AGE"), 86400) // Default 86400 seconds (24 hours)

	// Get the Authorization header or cookie
	var accessToken string
	accessToken, err := c.Cookie("access_token")
	if err != nil || accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Ensure the header contains "Bearer" token
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
				accessToken = tokenParts[1]
			}
		}
	}

	// If access token is missing or invalid, check the refresh token
	if accessToken == "" {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No valid tokens provided"})
			c.Abort()
			return
		}

		// Validate the refresh token
		userID, role, err := utils.ValidateAndExtractToken(refreshToken, "REFRESH")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			c.Abort()
			return
		}

		// Generate new access and refresh tokens
		newAccessToken, err := utils.GenerateToken(models_user.UserJwt{ID: userID, Role: role}, "ACCESS")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access tokens"})
			c.Abort()
			return
		}
		newRefreshToken, err := utils.GenerateToken(models_user.UserJwt{ID: userID, Role: role}, "REFRESH")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh tokens"})
			c.Abort()
			return
		}

		// Set the new tokens in cookies
		c.SetCookie("access_token", newAccessToken, int(accessTokenAge), "/", "", false, true)
		c.SetCookie("refresh_token", newRefreshToken, int(refreshTokenAge), "/", "", false, true)

		// Attach user information to the request context
		c.Set("userID", userID)
		c.Set("role", role)

		// Continue processing the request
		c.Next()
		return
	}

	// Validate and extract user information from the access token
	userID, role, err := utils.ValidateAndExtractToken(accessToken, "ACCESS")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		c.Abort()
		return
	}

	// Verify the user exists in the database
	var user models_user.User
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
