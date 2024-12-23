package postgres

import (
	"log"
	"menu-server/src/config/env"
	"menu-server/src/models" // Adjust to the actual path
	models_dine "menu-server/src/models/dine"
	models_restaurant "menu-server/src/models/restaurant"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database instance
var DB *gorm.DB

// Aliased models for brevity
type (
	Plan                = models.Plan
	User                = models.User
	Restaurant          = models.Restaurant
	Menu                = models.Menu
	MenuItem            = models.MenuItem
	MenuCategory        = models.MenuCategory
	MenuItemOption      = models.MenuItemOption
	RestaurantOrder     = models_restaurant.Order
	DinePayment         = models_dine.DinePayment
	Subscription        = models.Subscription
	RestaurantOrderItem = models_restaurant.OrderItem

	PlanFeature            = models.PlanFeature
	PlanFeatureAssociation = models.PlanFeatureAssociation
)

// InitDB initializes the PostgreSQL database connection and runs migrations.
func InitDB() {
	// Load the DATABASE_URL from environment variables
	dsn := env.PostgresDatabaseVar["DATABASE_URL"]
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Open the database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Ensure the UUID extension is available
	if err := ensureUUIDExtension(); err != nil {
		log.Fatalf("Failed to ensure UUID extension: %v", err)
	}
	// Drop table
	// DB.Migrator().DropTable(&Subscription{})
	// DB.Migrator().DropColumn(&models.Plan{}, "duration")
	// Run migrations for all models
	if err := migrateModels(); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	log.Println("Database initialized successfully")
}

// ensureUUIDExtension ensures the "uuid-ossp" PostgreSQL extension exists.
func ensureUUIDExtension() error {
	return DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
}

// migrateModels runs the database schema migrations.
func migrateModels() error {
	return DB.AutoMigrate(
		&User{},
		&Plan{},
		&PlanFeature{},
		&PlanFeatureAssociation{},
		&Restaurant{},
		&Menu{},
		&MenuCategory{},
		&MenuItem{},
		&MenuItemOption{},
		&RestaurantOrder{},
		&RestaurantOrderItem{},
		&DinePayment{},
		&Subscription{},
	)
}
