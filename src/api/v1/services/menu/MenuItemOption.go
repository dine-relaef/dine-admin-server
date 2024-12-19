package menu

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateMenuItem handles the creation of a new menu item
// @Summary Create a new menu item option
// @Description Create a new menu item option
// @Tags Menu
// @Accept json
// @Produce json
// @Param item body models.AddMenuItemOptionData true "Item data"
// @Param menu_item_id path string true "Menu Item ID"
// @Param category_id path string true "Category ID"
// @Param menu_id path string true "Menu ID"
// @Param restaurant_id path string true "Restaurant ID"
// @Router /api/v1/{restaurant_id}/menus/{menu_id}/categories/{category_id}/items/{menu_item_id}/options [post]
func CreateMenuItemOption(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	menuItemID := c.Param("menu_item_id")
	restaurantID := c.Param("restaurant_id")

	if role != "admin" {
		if err := postgres.DB.Where("restaurant_admin_id = ? AND id = ?", userID, restaurantID).First(&models.Restaurant{}).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create category for this menu"})
			return
		}
	}

	var menuItemOption models.MenuItemOption
	if err := c.ShouldBindJSON(&menuItemOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingItem models.MenuItem
	if err := postgres.DB.Where("menu_item_id = ? AND name = ?",menuItemID, menuItemOption.Name).First(&existingItem).Error; err == nil {
		// If an item with the same name and category exists
		c.JSON(http.StatusConflict, gin.H{"error": "Item with the same name already exists in this category"})
		return
	}

	menuItemUUID, err := uuid.FromString(menuItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID format"})
		return
	}
	menuItemOption.MenuItemID = menuItemUUID
	if err := postgres.DB.Create(&menuItemOption).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item Created successfully", " Item Option": menuItemOption})
}
