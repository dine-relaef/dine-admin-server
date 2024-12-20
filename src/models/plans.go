package models

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Plan struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(50);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Features     datatypes.JSON `gorm:"type:jsonb" json:"features"` // JSON for flexible feature storage
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	IsActive     bool           `gorm:"type:boolean;default:true" json:"is_active"`
	Duration     string         `gorm:"type:varchar(50);not null;default:'1M'" json:"duration"` // Valid values: 1M, 6M, 12M
	DurationDays int            `gorm:"type:int;not null" json:"duration_days"`                // Number of days for the plan
	TrialPeriod  bool           `gorm:"type:boolean;default:true" json:"trial_period"`
}

// BeforeSave is a GORM hook that validates the duration and calculates the number of days
func (p *Plan) BeforeSave(tx *gorm.DB) (err error) {
	validDurations := map[string]int{
		"1M": 30,   // 1 Month = 30 Days
		"6M": 180,  // 6 Months = 180 Days
		"12M": 365, // 12 Months = 365 Days
	}

	// Validate the duration
	durationDays, isValid := validDurations[p.Duration]
	if !isValid {
		return errors.New("invalid duration value, must be '1M', '6M', or '12M'")
	}

	// Set the calculated duration days
	p.DurationDays = durationDays
	return nil
}
