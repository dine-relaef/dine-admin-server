package routes

import (
	middleware "menu-server/src/api/v1/middleware"
	menu "menu-server/src/api/v1/services/menu"

	"github.com/gin-gonic/gin"
)

func SetupMenuRoutes(menuGroup *gin.RouterGroup) {
	// Routes for Menus
	menuGroup.POST("/", middleware.Authenticate, menu.CreateMenu) // Create a menu
	menuGroup.GET("/", menu.GetMenus)                             // Get all menus, supports ?restaurant_id=
	menuGroup.GET("/:menu_id", menu.GetMenuByID)                  // Get a specific menu by ID
	menuGroup.PUT("/:menu_id", menu.UpdateMenu)                   // Update a menu by ID
	menuGroup.DELETE("/:menu_id", menu.DeleteMenu)                // Delete a menu by ID

	// Nested Routes: Categories under a Menu
	categoriesGroup := menuGroup.Group("/:menu_id/categories")
	{
		categoriesGroup.POST("/", middleware.Authenticate, menu.CreateMenuCategory) // Create a category for a specific menu
		categoriesGroup.GET("/", menu.GetMenuCategories)                            // Get all categories for a specific menu
		categoriesGroup.GET("/:category_id", menu.GetMenuCategoryByID)              // Get a specific category by ID
		categoriesGroup.PUT("/:category_id", menu.UpdateMenuCategory)               // Update a category by ID
		categoriesGroup.DELETE("/:category_id", menu.DeleteMenuCategory)            // Delete a category by ID
	}

	// Nested Routes: Items under a Category
	itemsGroup := menuGroup.Group("/:menu_id/categories/:category_id/items")
	{
		itemsGroup.POST("/", middleware.Authenticate, menu.CreateMenuItem) // Create a menu item in a specific category
		itemsGroup.GET("/", menu.GetMenuItems)                             // Get all items for a specific category
		itemsGroup.GET("/:item_id", menu.GetMenuItemByID)                  // Get a specific item by ID
		itemsGroup.PUT("/:item_id", menu.UpdateMenuItem)                   // Update a menu item by ID
		itemsGroup.DELETE("/:item_id", menu.DeleteMenuItem)                // Delete a menu item by ID
	}

}
