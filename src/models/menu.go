package models

import (
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

type Menu struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null" json:"restaurant_id"`                       // Foreign key
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE;" json:"-"` // Back relationship

	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	MenuCategory []MenuCategory `gorm:"foreignKey:MenuID" json:"menu_category"` // One-to-Many relationship
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type MenuCategory struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MenuID      uuid.UUID  `gorm:"type:uuid;not null" json:"menu_id"`                       // Foreign key
	Menu        Menu       `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE;" json:"-"` // Back relationship
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	ImageURL    string     `gorm:"type:varchar(255)" json:"image_url"`
	MenuItems   []MenuItem `gorm:"foreignKey:MenuCategoryID" json:"menu_items"` // One-to-Many relationship
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type MenuItem struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MenuCategoryID uuid.UUID    `gorm:"type:uuid;not null" json:"menu_category_id"`                      // Foreign key
	MenuCategory   MenuCategory `gorm:"foreignKey:MenuCategoryID;constraint:OnDelete:CASCADE;" json:"-"` // Back relationship
	Name           string       `gorm:"type:varchar(100);not null" json:"name"`
	Description    string       `gorm:"type:text" json:"description"`
	Price          float64      `gorm:"type:decimal(10,2);not null" json:"price"`
	ImageURL       string       `gorm:"type:varchar(255)" json:"image_url"`
	IsVegetarian   bool         `gorm:"type:boolean;default:false" json:"is_vegetarian"`
	IsAvailable    bool         `gorm:"type:boolean;default:true" json:"is_available"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}
