package services

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
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
// @Param user body models.CreateUserData true "User data"
// @Router /api/v1/auth/register [post]
func RegisterUser(c *gin.Context) {
	var user models.User

	// Bind JSON to the User model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new UUID for the user
	newUUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.ID = newUUID

	// Save user to the database
	if err := postgres.DB.Create(&user).Error; err != nil {
		// Check for duplicate key violation
		if strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Already Exists"})
			return
		}
		// Handle other errors (optional)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}
