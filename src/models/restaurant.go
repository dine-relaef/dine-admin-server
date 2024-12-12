package models

import (
	"time"

	"github.com/gofrs/uuid"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	"gorm.io/gorm"
)

type Restaurant struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	PureVeg   bool      `gorm:"type:boolean;default:false" json:"pure_veg"`
	Location  string    `gorm:"type:varchar(255);not null" json:"location"`
	Address   string    `gorm:"type:varchar(255);not null" json:"address"`
	ImageURL  string    `gorm:"type:varchar(255)" json:"image_url"`
	AdminID   uuid.UUID `gorm:"type:uuid;not null" json:"admin_id"`
	Admin     User      `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE;" json:"-"`
	Menu      []Menu    `gorm:"foreignKey:RestaurantID" json:"menu"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID, err = uuid.NewV4()
	return err
}
