package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Location represents the geographical location of a restaurant.
type Location struct {
	Address   string  `json:"address"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (l *Location) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to scan Location: value is not []byte or string")
	}

	return json.Unmarshal(bytes, l)
}

// Restaurant represents the restaurant entity in the database.
type Restaurant struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	PureVeg  bool      `gorm:"type:boolean;default:false" json:"pure_veg"`
	Location Location  `gorm:"type:json" json:"location"`

	ImageURL       string    `gorm:"type:varchar(255)" json:"image_url"`
	AdminID        uuid.UUID `gorm:"type:uuid;not null" json:"admin_id"`
	BannerImageUrl string    `json:"banner_image_url"`
	LogoImageUrl   string    `json:"logo_image_url"`
	Admin          User      `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE;" json:"-"`
	Menu           []Menu    `gorm:"foreignKey:RestaurantID" json:"menu"`
	SubscriptionID uuid.UUID `gorm:"type:uuid" json:"subscription_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Description    string    `gorm:"type:varchar(500)" json:"description"`
	Phone          string    `gorm:"type:varchar(15);default" json:"phone"`
	Email          string    `gorm:"type:varchar(100);not null" json:"email"`
	IsActive       bool      `gorm:"type:boolean;default:true" json:"is_active"`
	HasParking     bool      `gorm:"type:boolean;default:false" json:"has_parking"`
	HasPickup      bool      `gorm:"type:boolean;default:false" json:"has_delivery"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID, err = uuid.NewV4()
	}
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
	IsActive       bool     `json:"is_active"`
	HasParking     bool     `json:"has_parking"`
	HasPickup      bool     `json:"has_delivery"`
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
	IsActive       bool      `json:"is_active"`
	HasParking     bool      `json:"has_parking"`
	HasPickup      bool      `json:"has_delivery"`
	NumberOfTables int       `json:"number_of_tables"`
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
	IsActive       bool     `json:"is_active"`
	HasParking     bool     `json:"has_parking"`
	HasPickup      bool     `json:"has_delivery"`
	NumberOfTables int      `json:"number_of_tables"`
}
