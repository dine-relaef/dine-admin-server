package services

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateUser handles the creation of a new user
func CreateUser(c *gin.Context) {
	var user models.User

	// Bind JSON to the User model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a user with the same phone number or email already exists
	var existingUser models.User
	if err := postgres.DB.Where("phone = ? OR email = ?", user.Phone, user.Email).First(&existingUser).Error; err == nil {
		// If a user is found, return a conflict error
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with the same phone number or email already exists",
		})
		return
	}

	// Generate a new UUID for the user
	newUUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}
	user.ID = newUUID

	// Save user to the database
	if err := postgres.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUsers retrieves all users from the postgres
func GetUsers(c *gin.Context) {
	var users []models.User

	if err := postgres.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users Found", "users": users})
}

// GetUserByID retrieves a specific user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := postgres.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser updates a user's information by ID
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := postgres.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

// DeleteUser removes a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := postgres.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
