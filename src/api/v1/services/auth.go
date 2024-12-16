package services

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	utils "menu-server/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// @BasePath /api/v1
// RegisterUser handles the creation of a new user
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.RegisterUserData true "User data"
// @Router /api/v1/auth/register [post]
func RegisterUser(c *gin.Context) {
	var userData models.RegisterUserData

	// Bind JSON to the User model
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Generate a new UUID for the user
	newUUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate user ID"})
		return
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user instance
	user := models.User{
		ID:       newUUID,
		Name:     userData.Name,
		Email:    userData.Email,
		Password: hashedPassword,
		Phone:    userData.Phone,
	}

	// Save user to the database
	if err := postgres.DB.Create(&user).Error; err != nil {
		// Check for duplicate key violation (e.g., unique email)
		if strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateToken(models.UserJwt{ID: user.ID.String(), Role: user.Role}, "ACCESS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateToken(models.UserJwt{ID: user.ID.String(), Role: user.Role}, "REFRESH")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Set tokens as cookies
	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
		"tokens": gin.H{
			"access":  accessToken,
			"refresh": refreshToken,
		},
	})
}

// LoginUser handles the login of a user
// @Summary Login a user
// @Description Login a user in the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginUserData true "User credentials"
// @Router /api/v1/auth/login [post]
func LoginUser(c *gin.Context) {
	var user models.User
	var loginData models.LoginUserData

	// Bind JSON to the LoginUserData model
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Check if user exists in the database
	if err := postgres.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify the password
	if !utils.CheckPassword(user.Password, loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateToken(models.UserJwt{ID: user.ID.String(), Role: user.Role}, "ACCESS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateToken(models.UserJwt{ID: user.ID.String(), Role: user.Role}, "REFRESH")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Set tokens as cookies
	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"tokens": gin.H{
			"access":  accessToken,
			"refresh": refreshToken,
		},
	})
}

// RefreshToken generates a new access token using a valid refresh token
// @Summary Refresh Access Token
// @Description Generate a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token header string true "Refresh Token"
// @Router /api/v1/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing or invalid"})
		return
	}

	// Validate and extract claims from the refresh token
	userID, role, err := utils.ValidateAndExtractToken(refreshToken, "REFRESH")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate a new access token
	accessToken, err := utils.GenerateToken(models.UserJwt{ID: userID, Role: role}, "ACCESS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Set the new access token as a cookie
	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Access token refreshed successfully",
		"access":  accessToken,
	})
}

// LogoutUser clears user authentication cookies
// @Summary Logout a user
// @Description Clear user authentication cookies
// @Tags Auth
// @Produce json
// @Router /api/v1/auth/logout [post]
func LogoutUser(c *gin.Context) {
	// Clear authentication cookies
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}