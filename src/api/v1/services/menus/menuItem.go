package services_menu

import (
	postgres "dine-server/src/config/database"
	models_menu "dine-server/src/models/menu"
	utils "dine-server/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateMenuItem handles the creation of a new menu item
// @Summary Create a new menu item
// @Description Create a new menu item
// @Tags Menu Category Item
// @Accept json
// @Produce json
// @Param item body models_menu.AddMenuItemData true "Item data"
// @Param category_id path string true "Category ID"
// @Param menu_id path string true "Menu ID"
// @Param restaurant_id path string true "Restaurant ID"
// @Router /api/v1/{restaurant_id}/menus/{menu_id}/categories/{category_id}/items [post]
func CreateMenuItem(c *gin.Context) {
	// Retrieve user info from context
	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create item for this category"})
		return
	}

	// Extract parameters
	categoryID := c.Param("category_id")
	menuID := c.Param("menu_id")

	// Parse and validate request body
	var addMenuItemData models_menu.AddMenuItemData
	if err := c.ShouldBindJSON(&addMenuItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if item already exists
	var existingItem models_menu.MenuItem
	if err := postgres.DB.Where("name = ? AND category_id = ?", addMenuItemData.Name, categoryID).First(&existingItem).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "An item with the same name already exists in this category"})
		return
	}

	// Parse UUIDs
	categoryUUID, err := uuid.FromString(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	menuUUID, err := uuid.FromString(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID format"})
		return
	}

	// Start a transaction
	tx := postgres.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Create new MenuItem
	newItemUUID := uuid.Must(uuid.NewV4())
	item := models_menu.MenuItem{
		ID:           newItemUUID,
		MenuID:       menuUUID,
		CategoryID:   categoryUUID,
		Name:         addMenuItemData.Name,
		Description:  &addMenuItemData.Description,
		ImageURL:     &addMenuItemData.ImageURL,
		IsVegetarian: addMenuItemData.IsVegetarian,
		IsAvailable:  addMenuItemData.IsAvailable,
	}

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item"})
		return
	}

	// Create MenuItemOptions if provided
	var createdOptions []models_menu.MenuItemOption
	for _, option := range addMenuItemData.ItemOptions {
		optionUUID := uuid.Must(uuid.NewV4())
		itemOption := models_menu.MenuItemOption{
			ID:         optionUUID,
			MenuItemID: newItemUUID,
			Name:       option.Name,
			Price:      option.Price,
		}
		if err := tx.Create(&itemOption).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item options"})
			return
		}
		createdOptions = append(createdOptions, itemOption)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Menu item created successfully",
		"item":    item,
		"options": createdOptions,
	})
}

// CreateMultipleMenuItems handles the creation of multiple menu items
// @Summary Create multiple menu items
// @Description Create multiple menu items
// @Tags Menu Category Item
// @Accept json
// @Produce json
// @Param items body []models_menu.AddMenuItemData true "Items data"
// @Param category_id path string true "Category ID"
// @Param menu_id path string true "Menu ID"
// @Param restaurant_id path string true "Restaurant ID"
// @Router /api/v1/{restaurant_id}/menus/{menu_id}/categories/{category_id}/items [post]
func CreateMultipleMenuItems(c *gin.Context) {
	// Retrieve user info from context
	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create item for this category"})
		return
	}

	// Extract parameters
	categoryID := c.Param("category_id")
	menuID := c.Param("menu_id")

	// Parse UUIDs
	categoryUUID, err := uuid.FromString(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	menuUUID, err := uuid.FromString(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID format"})
		return
	}
	// Parse and validate request body
	var addMenuItemData []models_menu.AddMenuItemData
	if err := c.ShouldBindJSON(&addMenuItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := postgres.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Create new MenuItem
	var createdItems []models_menu.MenuItem
	for _, itemData := range addMenuItemData {
		// Check if item already exists
		tx := postgres.DB.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Create new MenuItem
		newItemUUID := uuid.Must(uuid.NewV4())
		item := models_menu.MenuItem{
			ID:           newItemUUID,
			MenuID:       menuUUID,
			CategoryID:   categoryUUID,
			Name:         itemData.Name,
			Description:  &itemData.Description,
			ImageURL:     &itemData.ImageURL,
			IsVegetarian: itemData.IsVegetarian,
			IsAvailable:  itemData.IsAvailable,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item"})
			return
		}

		// Create MenuItemOptions if provided
		var createdOptions []models_menu.MenuItemOption
		for _, option := range itemData.ItemOptions {
			optionUUID := uuid.Must(uuid.NewV4())
			itemOption := models_menu.MenuItemOption{
				ID:         optionUUID,
				MenuItemID: newItemUUID,
				Name:       option.Name,
				Price:      option.Price,
			}
			if err := tx.Create(&itemOption).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item options"})
				return
			}
			createdOptions = append(createdOptions, itemOption)
		}

	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Menu items created successfully",
		"items":   createdItems,
	})
}

// GetMenuItems retrieves all items for a specific category
// @Summary Retrieve all items for a specific category
// @Description Retrieve all items for a specific category
// @Tags Menu
// @Produce json
// @Param category_id path string true "Category ID"
// @Router /api/v1/{restaurant_id}/menus/{menu_id}/categories/{category_id}/items [get]
func GetMenuItems(c *gin.Context) {
	categoryID := c.Query("category_id")

	var items []models_menu.MenuItem

	if err := postgres.DB.Where(" menu_category_id = ?", categoryID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items Found Successfully", "items": items})
}

// GetMenuItemByID retrieves a specific menu item by ID

func GetMenuItemByID(c *gin.Context) {
	itemID := c.Param("item_id")
	var item models_menu.MenuItem

	if err := postgres.DB.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item Got Successfully", "item": item})
}

// UpdateMenuItem updates a specific menu item by ID
func UpdateMenuItem(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update item for this category"})
		return
	}

	itemID := c.Param("item_id")
	var item models_menu.MenuItem

	if err := postgres.DB.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item Updated Successfully", "item": item})
}

// DeleteMenuItem deletes a specific menu item by ID
func DeleteMenuItem(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create menu for this restaurant"})
		return
	}

	itemID := c.Param("item_id")
	if err := postgres.DB.Delete(&models_menu.MenuItem{}, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
