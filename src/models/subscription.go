package models

import (
	"time"

	uuid "github.com/gofrs/uuid"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

// Subscription represents the subscription entity in the database.
type Subscription struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID uuid.UUID      `gorm:"type:uuid;not null" json:"restaurant_id"`
	PlanID       uuid.UUID      `gorm:"type:uuid;not null" json:"plan_id"`
	Plan         Plan           `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE;" json:"plan"`
	StartDate    time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time      `gorm:"type:date;not null" json:"end_date"`
	IsActive     bool           `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}