package models_user

import (
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	// "gorm.io/gorm"
	models_restaurant "dine-server/src/models/restaurants"
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID            uuid.UUID                      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name          string                         `gorm:"type:varchar(100);not null" json:"name"`
	Email         string                         `gorm:"type:varchar(100);not null;uniqueIndex:idx_email_phone" json:"email"`
	Phone         string                         `gorm:"type:varchar(20);not null;uniqueIndex:idx_email_phone" json:"phone"`
	VerifiedEmail bool                           `gorm:"type:boolean;default:false" json:"verified_email"`
	VerifiedPhone bool                           `gorm:"type:boolean;default:false" json:"verified_phone"`
	Password      string                         `gorm:"type:varchar(255);not null" json:"-"`
	Role          string                         `gorm:"type:varchar(50);not null;default:'restaurant_admin';check:role IN ('admin', 'restaurant_admin')" json:"role"`
	SignupSource  string                         `gorm:"type:varchar(50);not null;default:'website';check:signup_source IN ('website', 'google', 'facebook', 'apple')" json:"signup_source"`
	ProfileImage  string                         `gorm:"type:varchar(255)" json:"profile_image"`
	Restaurants   []models_restaurant.Restaurant `gorm:"foreignKey:AdminID;references:ID" json:"restaurants"`
	CreatedAt     time.Time                      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time                      `gorm:"autoUpdateTime" json:"updated_at"`
}

type UpdateUserDataByAdmin struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UpdateUserDataByUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(20)" json:"phone"`
	Role         string    `gorm:"type:varchar(50);not null;default:'restaurant_admin';check:role IN ('admin', 'restaurant_admin')" json:"role"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null" json:"restaurant_id"`
}

type UserJwt struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
