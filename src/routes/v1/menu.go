package routes

import (
	services "menu-server/src/api/v1/services"

	"github.com/gin-gonic/gin"
)

func SetupMenuRoutes(menuGroup *gin.RouterGroup) {
	// Routes for Menus
	menuGroup.POST("/", services.CreateMenu)           // Create a menu
	menuGroup.GET("/", services.GetMenus)              // Get all menus, supports ?restaurant_id=
	menuGroup.GET("/:menu_id", services.GetMenuByID)   // Get a specific menu by ID
	menuGroup.PUT("/:menu_id", services.UpdateMenu)    // Update a menu by ID
	menuGroup.DELETE("/:menu_id", services.DeleteMenu) // Delete a menu by ID

	// Nested Routes: Categories under a Menu
	categoriesGroup := menuGroup.Group("/:menu_id/categories")
	{
		categoriesGroup.POST("/", services.CreateMenuCategory)               // Create a category for a specific menu
		categoriesGroup.GET("/", services.GetMenuCategories)                 // Get all categories for a specific menu
		categoriesGroup.GET("/:category_id", services.GetMenuCategoryByID)   // Get a specific category by ID
		categoriesGroup.PUT("/:category_id", services.UpdateMenuCategory)    // Update a category by ID
		categoriesGroup.DELETE("/:category_id", services.DeleteMenuCategory) // Delete a category by ID
	}

	// Nested Routes: Items under a Category
	itemsGroup := categoriesGroup.Group("/:category_id/items")
	{
		itemsGroup.POST("/", services.CreateMenuItem)           // Create a menu item in a specific category
		itemsGroup.GET("/", services.GetMenuItems)              // Get all items for a specific category
		itemsGroup.GET("/:item_id", services.GetMenuItemByID)   // Get a specific item by ID
		itemsGroup.PUT("/:item_id", services.UpdateMenuItem)    // Update a menu item by ID
		itemsGroup.DELETE("/:item_id", services.DeleteMenuItem) // Delete a menu item by ID
	}
}
