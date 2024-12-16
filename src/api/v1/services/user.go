// src/api/v1/services/user.go
package services

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers retrieves all users or a specific user by ID if provided
// @Summary Get all users or a specific user by ID
// @Description Get all users in the system, or a specific user if the user_id query parameter is provided
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param user_id query string false "User ID to filter by"
// @Router /api/v1/users [get]
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
