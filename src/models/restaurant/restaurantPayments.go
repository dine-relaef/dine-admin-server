package models_restaurant

import (
	"time"

	"github.com/gofrs/uuid"
)

type RestaurantPayment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID  uuid.UUID `gorm:"type:uuid;not null" json:"restaurant_id"` // Foreign key
	OrderID       uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`      // Foreign key
	TransactionID *string   `gorm:"type:varchar(255)" json:"transaction_id,omitempty"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        *string   `gorm:"type:varchar(50);check:status IN ('successful','failed','pending');default:'pending'" json:"status,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

