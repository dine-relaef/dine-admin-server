package postgres

import (
	"dine-server/src/config/env" // Adjust to the actual path
	models_common "dine-server/src/models/Common"
	models_menu "dine-server/src/models/menu"
	models_order "dine-server/src/models/orders"
	models_payment "dine-server/src/models/payments"
	models_plan "dine-server/src/models/plans"
	models_promoCode "dine-server/src/models/promoCode"
	models_restaurant "dine-server/src/models/restaurants"
	models_subscription "dine-server/src/models/subscriptions"
	models_user "dine-server/src/models/users"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database instance
var DB *gorm.DB

// Aliased models for brevity
type (
	Plan                  = models_plan.Plan
	User                  = models_user.User
	Restaurant            = models_restaurant.Restaurant
	RestaurantBankAccount = models_restaurant.RestaurantBankAccount
	Menu                  = models_menu.Menu
	MenuItem              = models_menu.MenuItem
	MenuCategory          = models_menu.MenuCategory

	MenuItemOption  = models_menu.MenuItemOption
	RestaurantOrder = models_order.Order
	DinePayment     = models_payment.DinePayment

	DineOrder           = models_order.DineOrder
	Subscription        = models_subscription.Subscription
	RestaurantOrderItem = models_order.OrderItem

	PlanFeature            = models_plan.PlanFeature
	PlanFeatureAssociation = models_plan.PlanFeatureAssociation

	RestaurantsCount = models_common.RestaurantsCount

	DinePromoCode = models_promoCode.DinePromoCode
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
	// DB.Migrator().DropTable( &DinePromoCode{})
	// DB.Migrator().DropColumn(&models.DinePromoCode{}, "duration")
	// Run migrations for all models
	if err := migrateModels(); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Generate common tables

	log.Println("Database initialized successfully")
}

func GenrateCommonTable(table interface{}) {
	DB.Create(&table)
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
		&DineOrder{},
		&MenuItemOption{},
		&RestaurantOrder{},
		&RestaurantOrderItem{},
		&DinePayment{},
		&Subscription{},
		&RestaurantsCount{},
		&RestaurantBankAccount{},
		&DinePromoCode{},
	)
}
