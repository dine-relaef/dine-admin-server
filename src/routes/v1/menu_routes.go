package routes

import (
	controllers "menu-server/src/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetupMenuRoutes(router *gin.Engine) {
	MenuGroup := router.Group("/menus")
	{
		// Routes for Menus
		MenuGroup.POST("/", controllers.CreateMenu)         // Create a menu
		MenuGroup.GET("/", controllers.GetMenus)            // Get all menus, url + ?resturant_id=
		MenuGroup.GET("/:menu_id", controllers.GetMenuByID) // Get a specific menu by ID
		MenuGroup.PUT("/:menu_id", controllers.UpdateMenu)  // Update a menu by ID
		MenuGroup.DELETE("/:menu_id", controllers.DeleteMenu) // Delete a menu by ID
	}

		// Routes for Menu Categories under a Menu
		menuCategoryGroup := router.Group("/categories")
		{
			menuCategoryGroup.POST("/", controllers.CreateMenuCategory)             // Create a menu category for a specific menu
			menuCategoryGroup.GET("/", controllers.GetMenuCategories)              // Get all categories for a specific menu, url + ?menu_id=
			menuCategoryGroup.GET("/:category_id", controllers.GetMenuCategoryByID) // Get a specific category by ID
			menuCategoryGroup.PUT("/:category_id", controllers.UpdateMenuCategory) // Update a category by ID
			menuCategoryGroup.DELETE("/:category_id", controllers.DeleteMenuCategory) // Delete a category by ID
		}

		// Routes for Menu Items under a Category
		menuItemGroup := router.Group("/items")
		{
			menuItemGroup.POST("/", controllers.CreateMenuItem)             // Create a menu item in a specific category
			menuItemGroup.GET("/", controllers.GetMenuItems)               // Get all items for a specific category, url + ?category_id=
			menuItemGroup.GET("/:item_id", controllers.GetMenuItemByID)    // Get a specific item by ID
			menuItemGroup.PUT("/:item_id", controllers.UpdateMenuItem)     // Update a menu item by ID
			menuItemGroup.DELETE("/:item_id", controllers.DeleteMenuItem)  // Delete a menu item by ID
		}
	
}
