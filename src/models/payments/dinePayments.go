package models_payment

import (
	"time"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	"github.com/gofrs/uuid"
)

type DinePayment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID       uuid.UUID `gorm:"type:uuid;not null" json:"order_id"` // Foreign key
	Status        string    `gorm:"type:varchar(50);check:status IN ('successful','failed','pending');default:'pending';not null" json:"status"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	TransactionID string    `gorm:"type:varchar(255)" json:"transaction_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type RazerpayRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Receipt  string  `json:"receipt"`
}
