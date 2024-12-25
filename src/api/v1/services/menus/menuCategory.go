package services_menu

import (
	postgres "dine-server/src/config/database"
	models_menu "dine-server/src/models/menu"
	utils "dine-server/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateMenuCategory handles the creation of a new menu category
// @Summary Create a new category for a specific menu
// @Description Create a new category for a specific menu
// @Tags Menu Category
// @Accept json
// @Produce json
// @Param category body models_menu.AddMenuCategoryData true "Category data"
// @Param menu_id path string true "Menu ID"
// @Param restaurant_id path string true "Restaurant ID"
// @Router /api/v1/{restaurant_id}/menus/{menu_id}/categories [post]
func CreateMenuCategory(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create Category for this menu"})
		return
	}

	menuID := c.Param("menu_id")
	var category models_menu.MenuCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCategory models_menu.MenuCategory
	if err := postgres.DB.Where("name = ? AND menu_id = ?", category.Name, menuID).First(&existingCategory).Error; err == nil {
		// If a Category with the same name and Menu_id exists
		c.JSON(http.StatusConflict, gin.H{"error": "Category with the same name already exists in this Menu"})
		return
	}

	menuUUID, err := uuid.FromString(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID format"})
		return
	}

	category.MenuID = menuUUID
	if err := postgres.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category": category})
}

// GetMenuCategories retrieves all categories for a specific menu
func GetMenuCategories(c *gin.Context) {
	menuID := c.Query("menu_id")

	var categories []models_menu.MenuCategory
	if err := postgres.DB.Where("menu_id = ?", menuID).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categories Found Successfully", "categories": categories})
}

// GetMenuCategoryByID retrieves a specific menu category by ID
func GetMenuCategoryByID(c *gin.Context) {
	categoryID := c.Param("category_id")
	var category models_menu.MenuCategory

	if err := postgres.DB.First(&category, "id = ?", categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category Found Successfully", "category": category})
}

// UpdateMenuCategory updates a specific menu category by ID
func UpdateMenuCategory(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update category for this menu"})
		return
	}
	categoryID := c.Param("category_id")
	var category models_menu.MenuCategory

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

	c.JSON(http.StatusOK, gin.H{"message": "Category Updated Successfully", "category": category})
}

// DeleteMenuCategory deletes a specific menu category by ID
func DeleteMenuCategory(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")

	if isAdmin, err := utils.IsAuthorised(c, restaurantID); err != nil {
		// Unauthorized or Forbidden response
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete category for this menu"})
		return
	}

	categoryID := c.Param("category_id")
	if err := postgres.DB.Delete(&models_menu.MenuCategory{}, "id = ?", categoryID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
