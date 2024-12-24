package models_subscription

import (
	models_payment "menu-server/src/models/payments"
	models_plan "menu-server/src/models/plans"
	"time"

	uuid "github.com/gofrs/uuid"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

// Subscription represents the subscription entity in the database.
type Subscription struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	RestaurantID       uuid.UUID                  `gorm:"type:uuid;not null;ForeignKey:RestaurantID" json:"restaurant_id"`
	PlanID             uuid.UUID                  `gorm:"type:uuid;not null" json:"plan_id"`
	Plan               models_plan.Plan           `gorm:"foreignKey:PlanID;" json:"plan"`
	StartDate          time.Time                  `gorm:"type:date;not null" json:"start_date"`
	EndDate            time.Time                  `gorm:"type:date;not null" json:"end_date"`
	PaymentID          uuid.UUID                  `gorm:"type:uuid;not null" json:"payment_id"`
	Payment            models_payment.DinePayment `gorm:"foreignKey:PaymentID;" json:"-"`
	OrderID            uuid.UUID                  `gorm:"type:uuid;not null" json:"order_id"`
	AutoRenewal        bool                       `gorm:"type:boolean;default:false" json:"auto_renewal"`
	RenewalDate        time.Time                  `gorm:"type:date" json:"renewal_date"`
	LastRenewalDate    time.Time                  `gorm:"type:date" json:"last_renewal_date"`
	Canceled           bool                       `gorm:"type:boolean;default:false" json:"canceled"`
	CanceledAt         *time.Time                 `gorm:"type:timestamp" json:"canceled_at"`
	CancellationReason string                     `gorm:"type:varchar(255)" json:"cancellation_reason"`
	InGracePeriod      bool                       `gorm:"type:boolean;default:false" json:"in_grace_period"`
	GraceEndDate       *time.Time                 `gorm:"type:date" json:"grace_end_date"`
	CreatedAt          time.Time                  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time                  `gorm:"autoUpdateTime" json:"updated_at"`
}

type AddSubscriptionData struct {
	AutoRenewal bool `json:"auto_renewal" binding:"required"` // Whether auto-renewal
}
