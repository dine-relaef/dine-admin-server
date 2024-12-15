package services

import (
	"github.com/gin-gonic/gin"
	models "menu-server/src/models"
	postgres "menu-server/src/config/database"
	"net/http"
	"github.com/gofrs/uuid"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

// CreateMenu handles the creation of a new menu
func CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingMenu models.Menu
    if err := postgres.DB.Where("name = ? AND restaurant_id = ?", menu.Name, menu.RestaurantID).First(&existingMenu).Error; err == nil {
        // If an Category with the same name and Menu_id exists
        c.JSON(http.StatusConflict, gin.H{"error": "This Resturant has Menu with same name"})
        return
    }

	if err := postgres.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu"})
		return
	}

	

	c.JSON(http.StatusCreated, gin.H{"message" : "Menu Created Successfully", "menu": menu})
}

// GetMenus retrieves all menus
func GetMenus(c *gin.Context) {
	restaurantIdParam := c.Query("restaurant_id")
	restaurantID, err := uuid.FromString(restaurantIdParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
        return
    }
	
	var menus []models.Menu
	if err := postgres.DB.Where("restaurant_id = ?", restaurantID).Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Menus Found Successfully", "menus": menus})
}

// GetMenuByID retrieves a specific menu by ID
func GetMenuByID(c *gin.Context) {
	id := c.Param("menu_id")
	var menu models.Menu

	if err := postgres.DB.First(&menu, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Menu Got Successfully", "menu": menu})
}

// UpdateMenu updates a specific menu by ID
func UpdateMenu(c *gin.Context) {
	id := c.Param("menu_id")
	var menu models.Menu

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

	c.JSON(http.StatusOK, gin.H{"message" : "Menu Updated Successfully", "menu": menu})
}

// DeleteMenu deletes a specific menu by ID
func DeleteMenu(c *gin.Context) {
	id := c.Param("menu_id")
	if err := postgres.DB.Delete(&models.Menu{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}

// CreateMenuCategory handles the creation of a new menu category
func CreateMenuCategory(c *gin.Context) {
	var category models.MenuCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCategory models.MenuCategory
    if err := postgres.DB.Where("name = ? AND menu = ?", category.Name, category.MenuID).First(&existingCategory).Error; err == nil {
        // If an Category with the same name and Menu_id exists
        c.JSON(http.StatusConflict, gin.H{"error": "Category with the same name already exists in this Menu"})
        return
    }

	// category.MenuID = uuid(menuID)
	if err := postgres.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message" : "Category created successfully", "category": category})
}

// GetMenuCategories retrieves all categories for a specific menu
func GetMenuCategories(c *gin.Context) {
	menuID := c.Query("menu_id")

	var categories []models.MenuCategory
	if err := postgres.DB.Where("menu_id = ?", menuID).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Categories Found Successfully", "categories": categories})
}

// GetMenuCategoryByID retrieves a specific menu category by ID
func GetMenuCategoryByID(c *gin.Context) {
	categoryID := c.Param("category_id")
	var category models.MenuCategory

	if err := postgres.DB.First(&category, "id = ?", categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Category Found Successfully", "category": category})
}

// UpdateMenuCategory updates a specific menu category by ID
func UpdateMenuCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var category models.MenuCategory

	if err := postgres.DB.First(&category, "id = ?", categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Category Updated Successfully", "category": category})
}

// DeleteMenuCategory deletes a specific menu category by ID
func DeleteMenuCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	if err := postgres.DB.Delete(&models.MenuCategory{}, "id = ?", categoryID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// CreateMenuItem handles the creation of a new menu item
func CreateMenuItem(c *gin.Context) {
    var item models.MenuItem
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	var existingItem models.MenuItem

    if err := postgres.DB.Where("name = ? AND menu_category_id = ?", item.Name, item.MenuCategoryID).First(&existingItem).Error; err == nil {
        // If an item with the same name and category exists
        c.JSON(http.StatusConflict, gin.H{"error": "Item with the same name already exists in this category"})
        return
    }

    if err := postgres.DB.Create(&item).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message" : "Item Created successfully", "item": item})
}

// GetMenuItems retrieves all items for a specific category
func GetMenuItems(c *gin.Context) {
	categoryID := c.Query("category_id")

	var items []models.MenuItem

	if err := postgres.DB.Where(" menu_category_id = ?", categoryID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Items Found Successfully", "items": items})
}

// GetMenuItemByID retrieves a specific menu item by ID
func GetMenuItemByID(c *gin.Context) {
	itemID := c.Param("item_id")
	var item models.MenuItem

	if err := postgres.DB.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Item Got Successfully", "item": item})
}

// UpdateMenuItem updates a specific menu item by ID
func UpdateMenuItem(c *gin.Context) {
	itemID := c.Param("item_id")
	var item models.MenuItem

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

	c.JSON(http.StatusOK, gin.H{"message" : "Item Updated Successfully", "item": item})
}

// DeleteMenuItem deletes a specific menu item by ID
func DeleteMenuItem(c *gin.Context) {
	itemID := c.Param("item_id")
	if err := postgres.DB.Delete(&models.MenuItem{}, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
