package utils

import (
    "fmt"
    postgres "menu-server/src/config/database"
	models "menu-server/src/models"
    	"github.com/gin-gonic/gin"
)
func IsAuthorised(c *gin.Context, restaurantID string) (bool, error) {
    // Check if the role is admin
    role, _ := c.Get("role")
    if role == "admin" {
        return true, nil
    }

    // If role is not admin, check if user has access to the restaurant
    userID, exists := c.Get("userID")
    if !exists {
        return false, fmt.Errorf("unauthorized: user ID not found")
    }

    var restaurant models.Restaurant
    if err := postgres.DB.Where("id = ? AND restaurant_admin_id = ?", restaurantID, userID).First(&restaurant).Error; err != nil {
        return false, fmt.Errorf("forbidden: you are not allowed to access this restaurant")
    }

    return false, nil
}
