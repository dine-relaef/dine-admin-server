// src/api/v1/services/user.go
package services_user

import (
	postgres "menu-server/src/config/database"
	models_user "menu-server/src/models/users"
	utils "menu-server/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers retrieves all users
// @Summary Get all users
// @Description Get all users in the system
// @Tags Admin
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/users/get-all [get]
func GetAllUsers(c *gin.Context) {
	var users []models_user.User

	if err := postgres.DB.Preload("Restaurants").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users Found", "users": users})
}

// GetUserByID retrieves a specific user by ID
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags Admin
// @Produce json
// @Param id path string true "User ID"
// @Security ApiKeyAuth
// @Router /api/v1/users/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models_user.User

	if err := postgres.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser updates a user's information by ID
// @Summary Update a user by ID
// @Description Update a user by ID
// @Tags Admin
// @Produce json
// @Param id path string true "User ID"
// @Param user body models_user.UpdateUserDataByAdmin true "User data"
// @Security ApiKeyAuth
// @Router /api/v1/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models_user.User

	if err := postgres.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models_user.User
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
// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags Admin
// @Produce json
// @Param id path string true "User ID"
// @Security ApiKeyAuth
// @Router /api/v1/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if err := postgres.DB.Delete(&models_user.User{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUser retrieves the user information
// @Summary Get user information
// @Description Get the user information
// @Tags User
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/users [get]
func GetUser(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var user models_user.User

	if err := postgres.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser updates the user information
// @Summary Update user information
// @Description Update the user information
// @Tags User
// @Produce json
// @Param user body models_user.UpdateUserDataByUser true "User data"
// @Security ApiKeyAuth
// @Router /api/v1/users [put]
func UpdateUserByUser(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var user models_user.User

	if err := postgres.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models_user.UpdateUserDataByUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		input.Password = hashedPassword
	}

	if input.Password != "" {
		user.Password = input.Password
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Phone != "" {
		user.Phone = input.Phone
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
