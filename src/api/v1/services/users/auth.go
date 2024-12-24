package services_user

import (
	"context"
	"encoding/json"
	postgres "menu-server/src/config/database"
	"menu-server/src/config/env"
	models_user "menu-server/src/models/users"
	utils "menu-server/src/utils"
	"net/http"
	"net/url"
	"strings"
	"time"

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
// @Param user body models_user.RegisterUserData true "User data"
// @Router /api/v1/auth/register [post]
func RegisterUser(c *gin.Context) {
	var userData models_user.RegisterUserData

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
	user := models_user.User{
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
	accessToken, err := utils.GenerateToken(models_user.UserJwt{ID: user.ID.String(), Role: user.Role}, "ACCESS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateToken(models_user.UserJwt{ID: user.ID.String(), Role: user.Role}, "REFRESH")
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
// @Param credentials body models_user.LoginUserData true "User credentials"
// @Router /api/v1/auth/login [post]
func LoginUser(c *gin.Context) {
	var user models_user.User
	var loginData models_user.LoginUserData

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
	if user.SignupSource != "website" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login source"})
		return
	}

	// Verify the password
	if !utils.CheckPassword(user.Password, loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateToken(models_user.UserJwt{ID: user.ID.String(), Role: user.Role}, "ACCESS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateToken(models_user.UserJwt{ID: user.ID.String(), Role: user.Role}, "REFRESH")
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

// GoogleLogin initiates the Google OAuth2 flow
// @Summary Initiate Google OAuth2 login
// @Description Get Google OAuth2 URL with state parameter
// @Tags Auth
// @Accept json
// @Produce json
// @Router /api/v1/auth/google [get]
func GoogleLogin(c *gin.Context) {
	// Generate random state
	state := utils.GenerateState()

	// Store state in cookie
	c.SetCookie(
		"oauth_state",
		state,
		int(time.Now().Add(15*time.Minute).Unix()),
		"/",
		"",    // empty domain for testing
		false, // not secure for testing
		true,  // HttpOnly
	)

	url := env.Config.AuthCodeURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GoogleCallback handles the Google OAuth2 callback
// @Summary Google OAuth2 Callback
// @Description Handle Google OAuth2 callback
// @Tags Auth
// @Accept json
// @Param code query string true "OAuth2 authorization code"
// @Param state query string true "OAuth2 state parameter"
// @Produce json
// @Router /api/v1/auth/google/callback [post]
func GoogleCallback(c *gin.Context) {
	// Retrieve the authorization code from query parameters
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is missing"})
		return
	}

	// No need to decode the code parameter
	decodedCodeURL, err := url.QueryUnescape(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode code parameter"})
		return
	}
	// Validate the state parameter to prevent CSRF attacks
	state := c.Query("state")
	decodeStateURl, err := url.QueryUnescape(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode state parameter"})
		return
	}
	expectedState, _ := c.Cookie("oauth_state")
	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	if string(decodeStateURl) != expectedState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Exchange the authorization code for an access token
	token, err := env.Config.Exchange(context.Background(), string(decodedCodeURL))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the access token to fetch user info
	client := env.Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	// Decode the user info
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}
	var user models_user.User
	newAccount := false

	// Check if user exists in the database
	if err := postgres.DB.Where("email = ?", userInfo["email"].(string)).First(&user).Error; err != nil {
		// User not found, create a new user
		newAccount = true
		user = models_user.User{
			ID:            uuid.Must(uuid.NewV4()),
			Name:          userInfo["name"].(string),
			Email:         userInfo["email"].(string),
			VerifiedEmail: userInfo["verified_email"].(bool),
			ProfileImage:  userInfo["picture"].(string),
			SignupSource:  "google",
		}
		if err := postgres.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	accessToken, err := utils.GenerateToken(models_user.UserJwt{ID: user.ID.String(), Role: user.Role}, "ACCESS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateToken(models_user.UserJwt{ID: user.ID.String(), Role: user.Role}, "REFRESH")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Set tokens as cookies
	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600, "/", "", false, true)

	message := "Login successful"
	if newAccount {
		message = "Account created successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
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
	accessToken, err := utils.GenerateToken(models_user.UserJwt{ID: userID, Role: role}, "ACCESS")
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
