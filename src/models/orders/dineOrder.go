package models_order

import (
	"time"

	"github.com/gofrs/uuid"
)

type DineOrder struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID      uuid.UUID `gorm:"type:uuid;not null;" json:"restaurant_id"`
	RestaurantAdminID uuid.UUID `gorm:"type:uuid;not null" json:"restaurant_admin"`
	PlanID            uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`
	Amount            float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	PromoCode         string    `gorm:"type:varchar(50)" json:"discount_code"`
	DiscountAmount    float64   `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	Status            string    `gorm:"type:varchar(50);check:status IN ('successful','failed','pending');default:'pending';not null" json:"status"`
	Duration          string    `gorm:"type:varchar(50);not null" json:"type"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AddDineOrderData struct {
	RestaurantID uuid.UUID `json:"restaurant_id" binding:"required"`
	PlanID       uuid.UUID `json:"plan_id" binding:"required"`
	PromoCode    string    `json:"promo_code"`
	Duration     string    `json:"duration" binding:"required"`
}
