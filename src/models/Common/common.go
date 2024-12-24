package models_common

type RestaurantsCount struct {
	ID    int `gorm:"type:uuid;default:"1";primaryKey" json:"id"`
	Count int `gorm:"type:int;default:0" json:"count"`
}
