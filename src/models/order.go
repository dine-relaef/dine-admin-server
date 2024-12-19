package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Order represents the main order record
type Order struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RestaurantID  uuid.UUID   `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant    Restaurant  `gorm:"foreignKey:RestaurantID" json:"-"`
	CustomerEmail string      `gorm:"type:varchar(255);" json:"customer_email"`
	CustomerName  string      `gorm:"type:varchar(255);not null" json:"customer_name"`
	CustomerPhone string      `gorm:"type:varchar(255);not null" json:"customer_phone"`
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID" json:"items"` // Adjusted relationship
	PaymentType   string      `gorm:"type:varchar(20);check(payment_type in ('online', 'onsite'));not null" json:"payment_type"`
	Status        OrderStatus `gorm:"type:varchar(20);not null" json:"status"`
	OrderType     OrderType   `gorm:"type:varchar(20);not null" json:"order_type"`
	SubTotal      float64     `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Tax           float64     `gorm:"type:decimal(10,2);not null" json:"tax"`
	ServiceFee    float64     `gorm:"type:decimal(10,2);not null" json:"service_fee"`
	Total         float64     `gorm:"type:decimal(10,2);not null" json:"total"`
	Notes         *string     `gorm:"type:text" json:"notes"` // Optional field
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	CompletedAt   *time.Time  `json:"completed_at"`
}

// OrderItem represents individual items within an order
type OrderItem struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID        uuid.UUID      `gorm:"type:uuid;not null" json:"order_id"`
	MenuItemID     uuid.UUID      `gorm:"type:uuid;not null" json:"menu_item_id"`
	MenuItem       MenuItem       `gorm:"foreignKey:MenuItemID" json:"-"`
	MenuName       string         `gorm:"type:varchar(255)" json:"menu_name"`
	Quantity       int            `gorm:"type:int;not null" json:"quantity"`
	Price          float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Subtotal       float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	ItemOptionID   uuid.UUID      `gorm:"type:uuid" json:"item_option_id"`
	ItemOption     MenuItemOption `gorm:"foreignKey:ItemOptionID" json:"-"`
	ItemOptionName string         `gorm:"type:varchar(255)" json:"item_option_name"`
}

// OrderStatus represents the possible states of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusReady     OrderStatus = "READY"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeDineIn   OrderType = "DINEIN"
	OrderTypePickup   OrderType = "PICKUP" 
	OrderTypeDelivery OrderType = "DELIVERY"
)

// Data structs for creating orders
type CreateOrderItem struct {
	MenuItemID     uuid.UUID  `json:"menu_item_id" binding:"required"`
	Quantity       int        `json:"quantity" binding:"required,min=1"`
	ItemOptionID   *uuid.UUID `json:"item_option_id"`
	ItemOptionName *string    `json:"item_option_name"`
}

type CreateOrder struct {
	RestaurantID  uuid.UUID         `json:"restaurant_id" binding:"required"`
	CustomerEmail string            `json:"customer_email" binding:"email"`
	CustomerName  string            `json:"customer_name" binding:"required"`
	CustomerPhone string            `json:"customer_phone" binding:"required"`
	PaymentType   string            `json:"payment_type" binding:"required,oneof=online onsite"`
	OrderType     string            `json:"order_type" binding:"required,oneof=DINEIN PICKUP DELIVERY"`
	Notes         *string           `json:"notes"`
	Items         []CreateOrderItem `json:"items" binding:"required,dive"`
}
