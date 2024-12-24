package models_restaurant

import (
	"time"

	"github.com/gofrs/uuid"
)

type RestaurantBankAccount struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ContactID     string    `gorm:"type:varchar(100);not null" json:"contact_id"`
	FundAccID    string    `gorm:"type:varchar(100);not null" json:"fund_acc_id"`
	RestaurantID  uuid.UUID `gorm:"type:uuid;not null" json:"restaurant_id"`
	Email         string    `gorm:"type:varchar(100);not null" json:"email"`
	Phone         string    `gorm:"type:varchar(15);default" json:"phone"`
	BankName      string    `gorm:"type:varchar(100);not null" json:"bank_name"`
	AccountNumber string    `gorm:"type:varchar(100);not null" json:"account_number"`
	AccountName   string    `gorm:"type:varchar(100);not null" json:"account_holder"`
	IFSCCode      string    `gorm:"type:varchar(100);not null" json:"ifsc_code"`
	Branch        string    `gorm:"type:varchar(100);not null" json:"branch"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type AddBankAccount struct {
	BankName      string    `json:"bank_name"`
	RestaurantID  uuid.UUID `json:"restaurant_id"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_holder"`
	IFSCCode      string    `json:"ifsc_code"`
	Branch        string    `json:"branch"`
}
