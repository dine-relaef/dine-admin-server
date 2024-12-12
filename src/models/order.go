package models

import (
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	"gorm.io/datatypes"
)

type Order struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID uuid.UUID      `gorm:"type:uuid;not null" json:"restaurant_id"` // Foreign key
	CustomerID   uuid.UUID      `gorm:"type:uuid;not null" json:"customer_id"`   // Foreign key
	OrderItems   datatypes.JSON `gorm:"type:jsonb" json:"order_items"`           // JSON for items and quantities
	TotalPrice   float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status       string         `gorm:"type:varchar(50);check:status IN ('pending','preparing','ready','completed','cancelled');default:'pending'" json:"status"`
	PaymentType  string         `gorm:"type:varchar(50);check:payment_type IN ('online','onsite');default:'online';not null" json:"payment_type"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}
