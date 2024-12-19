package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Location represents the geographical location of a restaurant.
type Location struct {
	Address   string  ` json:"address"`
	Longitude float64 ` json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (l *Location) Scan(value interface{}) error {
	// If value is nil, set a default Location
	if value == nil {
		return nil
	}

	// Convert the value to []byte
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Location: value is not []byte")
	}

	// Unmarshal JSON into Location struct
	return json.Unmarshal(bytes, l)
}

// Value implements the driver.Valuer interface
func (l Location) Value() (driver.Value, error) {
	// Marshal Location struct to JSON
	return json.Marshal(l)
}

// Restaurant represents the restaurant entity in the database.
type Restaurant struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name              string    `gorm:"type:varchar(100);not null" json:"name"`
	Description       string    `gorm:"type:text" json:"description"`
	Phone             string    `gorm:"type:varchar(20);not null" json:"phone"`
	Email             string    `gorm:"type:varchar(100);not null" json:"email"`
	PureVeg           bool      `gorm:"type:boolean;default:false" json:"pure_veg"`
	Location          Location  `gorm:"type:jsonb;not null" json:"location"`
	BannerImageUrl    string    `gorm:"type:varchar(255)" json:"banner_image_url"`
	LogoImageUrl      string    `gorm:"type:varchar(255)" json:"logo_image_url"`
	RestaurantAdminID uuid.UUID `gorm:"type:uuid;not null" json:"restaurant_admin_id"`
	Admin             User      `gorm:"foreignKey:RestaurantAdminID;constraint:OnDelete:CASCADE;" json:"-"`
	Menu              []Menu    `gorm:"foreignKey:RestaurantID" json:"menu"`
	Orders            []Order   `gorm:"foreignKey:RestaurantID" json:"orders"`
	PlanID            uuid.UUID `gorm:"type:uuid;" json:"plan_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID, err = uuid.NewV4()
	return err
}

type AddRestaurantData struct {
	Name           string   `json:"name"`
	PureVeg        bool     `json:"pure_veg"`
	Description    string   `json:"description"`
	Location       Location `json:"location"`
	BannerImageUrl string   `json:"banner_image_url"`
	LogoImageUrl   string   `json:"logo_image_url"`
	Phone          string   `json:"phone"`
	Email          string   `json:"email"`
}

type ResponseRestaurantData struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	PureVeg        bool      `json:"pure_veg"`
	Description    string    `json:"description"`
	Location       Location  `json:"location"`
	BannerImageUrl string    `json:"banner_image_url"`
	LogoImageUrl   string    `json:"logo_image_url"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
}
type UpdateRestaurantData struct {
	Name           string   `json:"name"`
	PureVeg        bool     `json:"pure_veg"`
	Description    string   `json:"description"`
	Location       Location `json:"location"`
	BannerImageUrl string   `json:"banner_image_url"`
	LogoImageUrl   string   `json:"logo_image_url"`
	Phone          string   `json:"phone"`
	Email          string   `json:"email"`
}
