package services_menu

import (
	postgres "dine-server/src/config/database"
	models_menu "dine-server/src/models/menu"
	"net/http"

	utils "dine-server/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

// CreateMenu handles the creation of a new menu
// @Summary Create a new menu
// @Description Create a new menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param menu body models_menu.AddMenuData true "Menu data"
// @Param restaurant_id path string true "Restaurant ID"
// @Router /api/v1/{restaurant_id}/menus [post]
func CreateMenu(c *gin.Context) {

	var menu models_menu.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create menu for this restaurant"})
		return
	}

	var existingMenu models_menu.Menu
	if err := postgres.DB.Where("name = ? AND restaurant_id = ?", menu.Name, restaurantID).First(&existingMenu).Error; err == nil {
		// If a Menu with the same name and restaurant_id exists
		c.JSON(http.StatusConflict, gin.H{"error": "This Restaurant has a Menu with the same name"})
		return
	}

	RestaurantID, err := uuid.FromString(restaurantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID format"})
		return
	}
	menu.RestaurantID = RestaurantID
	if err := postgres.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Menu Created Successfully", "menu": menu})
}

// GetMenus retrieves all menus
// @Summary Retrieve all menus
// @Description Retrieve all menus
// @Tags Menu
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Router /api/v1/{restaurant_id}/menus [get]
func GetMenus(c *gin.Context) {
	restaurant_Id := c.Param("restaurant_id")

	var menus []models_menu.Menu
	if err := postgres.DB.Where("restaurant_id = ?", restaurant_Id).Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menus Found Successfully", "menus": menus})
}

// GetMenuByID retrieves a specific menu by ID
// @Summary Retrieve a specific menu by ID
// @Description Retrieve a specific menu by ID
// @Tags Menu
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param menu_id path string true "Menu ID"
// @Router /api/v1/{restaurant_id}/menus/{menu_id} [get]
func GetMenuByID(c *gin.Context) {
	menuID := c.Param("menu_id")
	restaurantID := c.Param("restaurant_id")

	// Fetch menu with related categories, items, and item options
	var menu models_menu.Menu
	if err := postgres.DB.
		Preload("Categories.MenuItems").
		Preload("Categories.MenuItems.ItemOptions"). // Preload ItemOptions for each MenuItem
		Where("id = ? AND restaurant_id = ?", menuID, restaurantID).
		First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Menu retrieved successfully",
		"menu":    menu,
	})
}

// UpdateMenu updates a specific menu by ID
func UpdateMenu(c *gin.Context) {
	id := c.Param("menu_id")
	var menu models_menu.Menu

	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update menu for this restaurant"})
		return
	}

	if err := postgres.DB.First(&menu, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu Updated Successfully", "menu": menu})
}

// DeleteMenu deletes a specific menu by ID
func DeleteMenu(c *gin.Context) {
	id := c.Param("menu_id")

	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create menu for this restaurant"})
		return
	}

	if err := postgres.DB.Delete(&models_menu.Menu{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}
