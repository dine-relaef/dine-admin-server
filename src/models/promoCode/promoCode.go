package models_promoCode

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

func (d *DinePromoCode) SetPlanIDs(ids []uuid.UUID) error {
	bytes, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	d.PlanIDs = bytes
	return nil
}

func (d *DinePromoCode) GetPlanIDs() ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := json.Unmarshal(d.PlanIDs, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

type DinePromoCode struct {
	ID           uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Code         string          `gorm:"type:varchar(255);unique;not null" json:"code"`
	Discount     float64         `gorm:"type:decimal(10,2);not null" json:"discount"`
	ValidFrom    string          `gorm:"type:varchar(255);not null" json:"valid_from"`
	ValidTo      string          `gorm:"type:varchar(255);not null" json:"valid_to"`
	MaxUses      int             `gorm:"type:int;not null" json:"max_uses"`
	DiscountType string          `gorm:"type:varchar(20);not null;check:discount_type IN ('amount','percentage')" json:"discount_type"`
	IsActive     bool            `gorm:"type:boolean;default:true" json:"is_active"`
	PlanIDs      json.RawMessage `gorm:"type:jsonb" json:"plan_ids"`
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type AddDinePromoCode struct {
	Code         string      `json:"code"`
	Discount     float64     `json:"discount"`
	Days         int         `json:"days"`
	MaxUses      int         `json:"max_uses"`
	PlanIDs      []uuid.UUID `json:"plan_ids"`
	IsActive     bool        `json:"is_active"`
	DiscountType string      `json:"discount_type"`
}

type RestaurantPromoCode struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID uuid.UUID   `gorm:"type:uuid;not null" json:"restaurant_id"` // Changed to UUID for consistency
	PromoCode    uuid.UUID   `gorm:"type:uuid;not null" json:"promo_code"`    // Changed to UUID for consistency
	ValidFrom    string      `gorm:"type:varchar(255);not null" json:"valid_from"`
	ValidTo      string      `gorm:"type:varchar(255);not null" json:"valid_to"`
	IsActive     bool        `gorm:"type:boolean;default:true" json:"is_active"`
	MenuItems    []uuid.UUID `gorm:"type:uuid[]"` // Array of UUIDs
	DiscountType string      `gorm:"type:varchar(20);not null;check:discount_type IN ('amount','percentage')" json:"discount_type"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
