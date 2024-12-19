package services

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	"menu-server/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// GetAllRestaurants godoc
// @Summary Retrieve all restaurants
// @Description Retrieve all restaurants
// @Tags Restaurant
// @Produce json
// @Router /api/v1/restaurants/get-all [get]
func GetAllRestaurants(c *gin.Context) {
	var restaurant []models.Restaurant

	if err := postgres.DB.Find(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurants"})
		return
	}

	if len(restaurant) == 0 {
		c.JSON(http.StatusOK, gin.H{"restaurants": []models.ResponseRestaurantData{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"restaurants": utils.RestaurantResponse(restaurant)})
}

// GetRestaurants godoc
// @Summary Retrieve all restaurants by User
// @Description Retrieve all restaurants by User
// @Tags Restaurant
// @Produce json
// @Router /api/v1/restaurants [get]
func GetRestaurants(c *gin.Context) {
	var restaurant []models.Restaurant
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if err := postgres.DB.Where("restaurant_admin_id = ?", userId).Find(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"restaurants": restaurant})
}

// GetRestaurantByID godoc
// @Summary Retrieve a restaurant by ID
// @Description Retrieve a restaurant by ID
// @Tags Restaurant
// @Produce json
// @Param id path string true "Restaurant ID"
// @Router /api/v1/restaurants/{id} [get]
func GetRestaurantByID(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	role, _ := c.Get("role")
	if role == "restaurant_admin" {
		if err := postgres.DB.Where("id = ? AND restaurant_admin_id = ?", id, userID).Find(&restaurant).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant ID not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"restaurant": restaurant})
		return
	}
	if err := postgres.DB.Where("id = ? ", id).Find(&restaurant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant ID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"restaurant": restaurant})
}

// CreateRestaurant godoc
// @Summary Create a new restaurant
// @Description Create a new restaurant
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param restaurant body models.AddRestaurantData true "Restaurant data"
// @Router /api/v1/restaurants [post]
func CreateRestaurant(c *gin.Context) {
	var restaurantData models.AddRestaurantData
	if err := c.ShouldBindJSON(&restaurantData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurntAdminIDStr, _ := c.Get("userID")
	restaurntAdminID, ok := restaurntAdminIDStr.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	restaurantAdminID, err := uuid.FromString(restaurntAdminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}
	var restaurant models.Restaurant = models.Restaurant{
		RestaurantAdminID: restaurantAdminID,
		Name:              restaurantData.Name,
		Description:       restaurantData.Description,
		Location:          restaurantData.Location,
		Phone:             restaurantData.Phone,
		PureVeg:           restaurantData.PureVeg,
		Email:             restaurantData.Email,
		LogoImageUrl:      restaurantData.LogoImageUrl,
		BannerImageUrl:    restaurantData.BannerImageUrl,
	}
	if err := postgres.DB.Create(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Restaurant created successfully",
		"restaurant": restaurant,
	})
}

// UpdateRestaurant godoc
// @Summary Update a restaurant
// @Description Update a restaurant by ID
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Param restaurant body models.UpdateRestaurantData true "Restaurant data"
// @Router /api/v1/restaurants/{id} [put]
func UpdateRestaurant(c *gin.Context) {
	
	id := c.Param("id")
	var restaurant models.Restaurant

	if err := postgres.DB.First(&restaurant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant ID not found"})
		return
	}

	role, _ := c.Get("role")
	if role != "admin" && role != "restaurant_admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role == "restaurant_admin" {
		userId, _ := c.Get("userID")
		if restaurant.RestaurantAdminID.String() != userId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	}

	var input models.Restaurant
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		restaurant.Name = input.Name
	}
	if input.Description != "" {
		restaurant.Description = input.Description
	}
	if input.PureVeg != restaurant.PureVeg {
		restaurant.PureVeg = input.PureVeg
	}
	if input.Phone != "" {
		restaurant.Phone = input.Phone
	}
	if input.Email != "" {
		restaurant.Email = input.Email
	}
	if input.BannerImageUrl != "" {
		restaurant.BannerImageUrl = input.BannerImageUrl
	}
	if input.LogoImageUrl != "" {
		restaurant.LogoImageUrl = input.LogoImageUrl
	}

	if err := postgres.DB.Save(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Restaurant updated successfully",
		"restaurant": restaurant,
	})
}

// DeleteRestaurant godoc
// @Summary Delete a restaurant
// @Description Delete a restaurant by ID
// @Tags Restaurant
// @Produce json
// @Param id path string true "Restaurant ID"
// @Router /api/v1/restaurants/{id} [delete]
func DeleteRestaurant(c *gin.Context) {

	id := c.Param("id")
	role, _ := c.Get("role")

	var restaurant models.Restaurant
	query := postgres.DB.Where("id = ?", id)
	if role == "restaurant_admin" {
		userId, _ := c.Get("userID")
		query = query.Where("restaurant_admin_id = ?", userId)
	}
	if err := query.Delete(&restaurant).Error; err != nil {
		if role == "restaurant_admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid restaurant ID"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
}
