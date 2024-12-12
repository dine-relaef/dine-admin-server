package models

import (
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

type Subscription struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null" json:"restaurant_id"` // Foreign key
	PlanID       uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`       // Foreign key
	StartDate    time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time `gorm:"type:date;not null" json:"end_date"`
	IsActive     bool      `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
