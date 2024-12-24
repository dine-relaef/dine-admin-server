package routes_v1

import (
	middleware "menu-server/src/api/v1/middleware"
	services_menu "menu-server/src/api/v1/services/menus"

	"github.com/gin-gonic/gin"
)

func SetupMenuRoutes(menuGroup *gin.RouterGroup) {
	// Routes for Menus
	menuGroup.POST("/", middleware.Authenticate, services_menu.CreateMenu) // Create a menu
	menuGroup.GET("/", services_menu.GetMenus)                             // Get all menus, supports ?restaurant_id=
	menuGroup.GET("/:menu_id", services_menu.GetMenuByID)                  // Get a specific menu by ID
	menuGroup.PUT("/:menu_id", services_menu.UpdateMenu)                   // Update a menu by ID
	menuGroup.DELETE("/:menu_id", services_menu.DeleteMenu)                // Delete a menu by ID

	// Nested Routes: Categories under a Menu
	categoriesGroup := menuGroup.Group("/:menu_id/categories")
	{
		categoriesGroup.POST("/", middleware.Authenticate, services_menu.CreateMenuCategory) // Create a category for a specific menu
		categoriesGroup.GET("/", services_menu.GetMenuCategories)                            // Get all categories for a specific menu
		categoriesGroup.GET("/:category_id", services_menu.GetMenuCategoryByID)              // Get a specific category by ID
		categoriesGroup.PUT("/:category_id", services_menu.UpdateMenuCategory)               // Update a category by ID
		categoriesGroup.DELETE("/:category_id", services_menu.DeleteMenuCategory)            // Delete a category by ID
	}

	// Nested Routes: Items under a Category
	itemsGroup := menuGroup.Group("/:menu_id/categories/:category_id/items")
	{
		itemsGroup.POST("/", middleware.Authenticate, services_menu.CreateMenuItem) // Create a menu item in a specific category
		itemsGroup.GET("/", services_menu.GetMenuItems)                             // Get all items for a specific category
		itemsGroup.GET("/:item_id", services_menu.GetMenuItemByID)                  // Get a specific item by ID
		itemsGroup.PUT("/:item_id", services_menu.UpdateMenuItem)                   // Update a menu item by ID
		itemsGroup.DELETE("/:item_id", services_menu.DeleteMenuItem)                // Delete a menu item by ID
	}

}
