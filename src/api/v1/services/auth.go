package services

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	"menu-server/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// @BasePath /api/v1
// CreateUser handles the creation of a new user
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new UUID for the user
	newUUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var user models.User = models.User{
		ID:       newUUID,
		Name:     userData.Name,
		Email:    userData.Email,
		Password: hashedPassword,
		Phone:    userData.Phone,
	}

	// Save user to the database
	if err := postgres.DB.Create(&user).Error; err != nil {
		// Check for duplicate key violation
		if strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Already Exists"})
			return
		}
		c.SetCookie("assess_token", "test", 3600, "/", "*", false, true)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	if err := postgres.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPassword(user.Password, loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.SetCookie("assess_token", "login", 3600, "/", "*", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}
