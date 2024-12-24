package models_plan

import (
	"time"

	"github.com/gofrs/uuid"
)

type Plan struct {
	ID                      uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name                    string                   `gorm:"type:varchar(50);not null" json:"name"`
	Description             string                   `gorm:"type:text" json:"description"`
	Price                   float64                  `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive                bool                     `gorm:"type:boolean;default:true" json:"is_active"`
	TrialPeriod             bool                     `gorm:"type:boolean;default:true" json:"trial_period"`
	PlanFeatureAssociations []PlanFeatureAssociation `gorm:"foreignKey:PlanID" json:"-"` // Corrected the field name
	CreatedAt               time.Time                `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time                `gorm:"autoUpdateTime" json:"updated_at"`
}

type PlanFeature struct {
	ID                      uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name                    string                   `gorm:"type:varchar(50);not null" json:"name"`
	Description             string                   `gorm:"type:text" json:"description"`
	PlanFeatureAssociations []PlanFeatureAssociation `gorm:"foreignKey:FeatureID" json:"-"` // Corrected the field name
	CreatedAt               time.Time                `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time                `gorm:"autoUpdateTime" json:"updated_at"`
}

type PlanFeatureAssociation struct {
	PlanID    uuid.UUID   `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_plan_feature" json:"plan_id"`
	FeatureID uuid.UUID   `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_plan_feature" json:"feature_id"`
	Plan      Plan        `gorm:"foreignKey:PlanID" json:"plan"`
	Feature   PlanFeature `gorm:"foreignKey:FeatureID" json:"feature"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type PlanResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	IsActive    bool          `json:"is_active"`
	TrialPeriod bool          `json:"trial_period"`
	Feature     []PlanFeature `json:"features"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type AddPlanData struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	IsActive    bool    `json:"is_active" binding:"default=false"`
	TrialPeriod bool    `json:"trial_period" binding:"required,default=false"`
}

type UpdatePlanData struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsActive    bool    `json:"is_active"`
	TrialPeriod bool    `json:"trial_period"`
}

type AddPlanFeatureData struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePlanFeatureData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PlanFeatureResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Plans       []Plan    `json:"plans"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
