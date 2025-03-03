package services

import (
	postgres "dine-server/src/config/database"
	models_common "dine-server/src/models/Common"
	models_restaurant "dine-server/src/models/restaurants"
	"dine-server/src/utils"
	"fmt"
	"log"
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
	var restaurant []models_restaurant.Restaurant

	if err := postgres.DB.Find(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurants"})
		return
	}

	if len(restaurant) == 0 {
		c.JSON(http.StatusOK, gin.H{"restaurants": []models_restaurant.ResponseRestaurantData{}})
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
	var restaurantDatas []models_restaurant.Restaurant

	// Retrieve user ID from context
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch restaurants associated with the user ID
	if err := postgres.DB.Where("admin_id = ?", userId).Find(&restaurantDatas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurants"})
		return
	}

	// Map fetched data to ResponseRestaurant structs
	var response []models_restaurant.ResponseRestaurantData
	for _, restaurant := range restaurantDatas {
		response = append(response, models_restaurant.ResponseRestaurantData{
			ID:             restaurant.ID,
			Name:           restaurant.Name,
			Description:    restaurant.Description,
			PureVeg:        restaurant.PureVeg,
			Location:       restaurant.Location,
			BannerImageUrl: restaurant.BannerImageUrl,
			LogoImageUrl:   restaurant.LogoImageUrl,
			Phone:          restaurant.Phone,
			Email:          restaurant.Email,
			IsActive:       restaurant.IsActive,
			HasParking:     restaurant.HasParking,
			HasPickup:      restaurant.HasPickup,
		})
	}

	c.JSON(http.StatusOK, gin.H{"restaurants": response})
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
	var restaurantData models_restaurant.Restaurant

	// Extract userID and role from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	role, _ := c.Get("role")

	// Restaurant Admin access check
	if role == "restaurant_admin" {
		if err := postgres.DB.Where("id = ? AND admin_id = ?", id, userID).First(&restaurantData).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
	} else {
		// General access
		if err := postgres.DB.Where("id = ?", id).First(&restaurantData).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant ID not found"})
			return
		}
	}

	// Map data to ResponseRestaurant struct
	response := models_restaurant.ResponseRestaurantData{
		ID:             restaurantData.ID,
		Name:           restaurantData.Name,
		Description:    restaurantData.Description,
		PureVeg:        restaurantData.PureVeg,
		Location:       restaurantData.Location,
		BannerImageUrl: restaurantData.BannerImageUrl,
		LogoImageUrl:   restaurantData.LogoImageUrl,
		Phone:          restaurantData.Phone,
		Email:          restaurantData.Email,
		IsActive:       restaurantData.IsActive,
		HasParking:     restaurantData.HasParking,
		HasPickup:      restaurantData.HasPickup,
	}

	c.JSON(http.StatusOK, gin.H{"restaurant": response})
}

// CreateRestaurant godoc
// @Summary Create a new restaurant
// @Description Create a new restaurant
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param restaurant body models_restaurant.AddRestaurantData true "Restaurant data"
// @Router /api/v1/restaurants [post]
func CreateRestaurant(c *gin.Context) {

	restaurntAdminIDStr, _ := c.Get("userID")

	var restaurantData models_restaurant.AddRestaurantData
	if err := c.ShouldBindJSON(&restaurantData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error 1": err.Error()})
		return
	}

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
	var restaurants []models_common.RestaurantsCount
	if err := postgres.DB.Find(&restaurants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(restaurants) <= 0 {
		var restaurant = models_common.RestaurantsCount{
			Count: 1,
		}
		if err := postgres.DB.Create(&restaurant).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		log.Println(restaurants[0].Count)
		var restaurant = models_common.RestaurantsCount{
			Count: restaurants[0].Count + 1,
		}
		if err := postgres.DB.Save(&restaurant).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	RestaurantCode := "R" + restaurantData.Location.StateCode + fmt.Sprintf("%04d", restaurants[0].Count+1)

	var restaurant = models_restaurant.Restaurant{
		AdminID:        restaurantAdminID,
		Name:           restaurantData.Name,
		Description:    restaurantData.Description,
		Location:       restaurantData.Location,
		RestaurantCode: RestaurantCode,
		Phone:          restaurantData.Phone,
		PureVeg:        restaurantData.PureVeg,
		Email:          restaurantData.Email,
		IsActive:       restaurantData.HasParking,
		HasParking:     restaurantData.HasParking,
		HasPickup:      restaurantData.HasPickup,
		LogoImageUrl:   restaurantData.LogoImageUrl,
		BannerImageUrl: restaurantData.BannerImageUrl,
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
// @Param restaurant body models_restaurant.UpdateRestaurantData true "Restaurant data"
// @Router /api/v1/restaurants/{id} [put]
func UpdateRestaurant(c *gin.Context) {

	id := c.Param("id")
	var restaurant models_restaurant.Restaurant

	// Fetch the restaurant by ID
	if err := postgres.DB.First(&restaurant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant ID not found"})
		return
	}

	// Role-based authorization
	role, _ := c.Get("role")
	if role != "admin" && role != "restaurant_admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the restaurant_admin is authorized for this restaurant
	if role == "restaurant_admin" {
		userID, _ := c.Get("userID")
		if restaurant.AdminID.String() != userID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	}

	// Bind request data to UpdateRestaurantData schema
	var resturantData models_restaurant.UpdateRestaurantData
	if err := c.ShouldBindJSON(&resturantData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dynamically map fields from UpdateRestaurantData to Restaurant
	updates := models_restaurant.Restaurant{
		Name:           resturantData.Name,
		Description:    resturantData.Description,
		PureVeg:        resturantData.PureVeg,
		Phone:          resturantData.Phone,
		Email:          resturantData.Email,
		BannerImageUrl: resturantData.BannerImageUrl,
		LogoImageUrl:   resturantData.LogoImageUrl,
		IsActive:       resturantData.IsActive,
		HasParking:     resturantData.HasParking,
		HasPickup:      resturantData.HasPickup,
	}

	// Perform updates only on the fields that are non-zero values
	if err := postgres.DB.Model(&restaurant).Updates(updates).Error; err != nil {
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

	var restaurant models_restaurant.Restaurant
	query := postgres.DB.Where("id = ?", id)
	if role == "restaurant_admin" {
		userId, _ := c.Get("userID")
		query = query.Where("admin_id = ?", userId)
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
