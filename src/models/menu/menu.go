package models_menu

import (
	"time"

	"github.com/gofrs/uuid"
)

// Menu Table
type Menu struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID uuid.UUID      `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Categories   []MenuCategory `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE;" json:"categories"`
	MenuItems    []MenuItem     `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE;" json:"items"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type MenuCategory struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MenuID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"menu_id"`
	Menu        Menu       `gorm:"foreignKey:MenuID" json:"-"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	ImageURL    *string    `gorm:"type:varchar(255)" json:"image_url"`
	Description *string    `gorm:"type:text" json:"description"`
	MenuItems   []MenuItem `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;" json:"menu_items"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type MenuItem struct {
	ID           uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MenuID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"menu_id"`
	Menu         Menu             `gorm:"foreignKey:MenuID" json:"-"`
	CategoryID   uuid.UUID        `gorm:"type:uuid;not null;index" json:"category_id"`
	Category     MenuCategory     `gorm:"foreignKey:CategoryID" json:"-"`
	Name         string           `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Description  *string          `gorm:"type:text" json:"description"`
	ImageURL     *string          `gorm:"type:varchar(255)" json:"image_url"`
	IsVegetarian bool             `gorm:"type:boolean;default:false" json:"is_vegetarian"`
	IsAvailable  bool             `gorm:"type:boolean;default:true" json:"is_available"`
	ItemOptions  []MenuItemOption `gorm:"foreignKey:MenuItemID;constraint:OnDelete:CASCADE;" json:"options"`
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

type MenuItemOption struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MenuItemID uuid.UUID `gorm:"type:uuid;not null;index" json:"menu_item_id"`
	MenuItem   MenuItem  `gorm:"foreignKey:MenuItemID" json:"-"`
	Name       string    `gorm:"type:varchar(50);not null" json:"name" validate:"required,min=1,max=50"`
	Price      float64   `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gt=0"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type AddMenuData struct {
	Name string `json:"name" binding:"required"`
}
type AddMenuCategoryData struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// MenuItemOption Table

type AddMenuItemData struct {
	Name         string                  `json:"name" binding:"required"`
	Description  string                  `json:"description"`
	ImageURL     string                  `json:"image_url"`
	IsVegetarian bool                    `json:"is_vegetarian"`
	IsAvailable  bool                    `json:"is_available"`
	ItemOptions  []AddMenuItemOptionData `json:"options"`
}

type AddMenuItemOptionData struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}
