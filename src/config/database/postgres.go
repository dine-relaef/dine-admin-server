package postgres

import (
	"log"
	models "menu-server/src/models" // replace with the actual path to your users package
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Plan models.Plan
type User models.User
type Restaurant models.Restaurant
type Menu models.Menu
type MenuItem models.MenuItem
type MenuCategory models.MenuCategory
type Order models.Order
type Payment models.Payment
type Subscription models.Subscription

func InitDB() {
	var err error

	// Read the DATABASE_URL environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Connect to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create UUID extension
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Printf("Error creating uuid-ossp extension: %v", err)
		// Note: You might want to handle this error more robustly
	}

	// Auto-migrate models
	err = DB.AutoMigrate(
		&User{},
		&Plan{},
		&Restaurant{},
		&Menu{},
		&MenuCategory{},
		&MenuItem{},
		&Order{},
		&Payment{},
		&Subscription{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}
