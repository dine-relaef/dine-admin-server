package models

type DinePromoCode struct {
	ID           string  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Code         string  `gorm:"type:varchar(255);not null" json:"code"`
	Discount     float64 `gorm:"type:decimal(10,2);not null" json:"discount"`
	ValidFrom    string  `gorm:"type:date;not null" json:"valid_from"`
	ValidTo      string  `gorm:"type:date;not null" json:"valid_to"`
	MaxUses      int     `gorm:"type:int;not null" json:"max_uses"`
	DiscountType string  `gorm:"type:varchar(20);check(discount_type in ('percentage', 'amount'));not null" json:"discount_type"`
	IsActive     bool    `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt    string  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    string  `gorm:"autoUpdateTime" json:"updated_at"`
}

type RestaurantPromoCode struct {
	ID           string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID string `gorm:"type:uuid;not null" json:"restaurant_id"`
	PromoCode    string `gorm:"type:uuid;not null" json:"promo_code"`
	ValidFrom    string `gorm:"type:date;not null" json:"valid_from"`
	ValidTo      string `gorm:"type:date;not null" json:"valid_to"`
	IsActive     bool   `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt    string `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    string `gorm:"autoUpdateTime" json:"updated_at"`
}

type AddPromoCodeData struct {
	Code         string  `json:"code"`
	Discount     float64 `json:"discount"`
	ValidFrom    string  `json:"valid_from"`
	ValidTo      string  `json:"valid_to"`
	MaxUses      int     `json:"max_uses"`
	DiscountType string  `json:"discount_type"`
}
