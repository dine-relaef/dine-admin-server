package models

import (
	"time"

	"github.com/gofrs/uuid"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	// "gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string       `gorm:"type:varchar(100);not null" json:"name"`
	Email       string       `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password    string       `gorm:"type:varchar(255);not null" json:"-"`
	Role        string       `gorm:"type:varchar(50);not null;default:'restaurant_admin';check:role IN ('admin', 'restaurant_admin')" json:"role"`
	Phone       string       `gorm:"type:varchar(20);unique;not null" json:"phone"`
	Restaurants []Restaurant `gorm:"foreignKey:AdminID" json:"restaurants"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

type UpdateUserData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type UserResponse struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(20)" json:"phone"`
	Role         string    `gorm:"type:varchar(50);not null;default:'restaurant_admin';check:role IN ('admin', 'restaurant_admin')" json:"role"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null" json:"restaurant_id"`
}
