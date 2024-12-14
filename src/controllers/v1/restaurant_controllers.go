package controllers

import (
	models "menu-server/src/models"
	postgres "menu-server/src/config/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Createrestaurant handles creating a new restaurant
func CreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save Restaurant to postgres
	if err := postgres.DB.Create(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}


	c.JSON(http.StatusCreated, gin.H{
		"message": "restaurants Paln created successfully",
		"restaurant": restaurant,
	})
}

// GetAllrestaurants retrieves all restaurants
func GetAllRestaurants(c *gin.Context) {
	var restaurant []models.Restaurant

	if err := postgres.DB.Find(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"restaurants": restaurant})
}

// GetrestaurantByID retrieves a restaurant by its ID
func GetRestaurantByID(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant

	if err := postgres.DB.First(&restaurant, "id = ?", id).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant ID not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"restaurant": restaurant})
}

// Updaterestaurant updates a restaurant by ID
func UpdateRestaurant(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant

	if err := postgres.DB.First(&restaurant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant ID not found"})
		return
	}

	var input models.Restaurant
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&restaurant).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant Plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "restaurant updated successfully",
		"restaurant": restaurant,
	})
}

// Deleterestaurant deletes a restaurant by ID
func DeleteRestaurant(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant
	if err := postgres.DB.Delete(&restaurant, "id = ?", id).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid restaurant ID"})
	}


	c.JSON(http.StatusOK, gin.H{"message" : "restaurant deleted successfully"})
}
