package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type EstateType string

const (
	Garage      EstateType = "garage"
	StoreRoom   EstateType = "store_room"
	Parking     EstateType = "parking_place"
	FlatDaily   EstateType = "flat_daily"
	FlatYearly  EstateType = "flat_yearly"
	ForBusiness EstateType = "for_business"
	Car         EstateType = "car"
	Other       EstateType = "other"
)

type Estate struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name" gorm:"not null" validate:"required"`
	Address      string        `json:"address" gorm:"not null" validate:"required"`
	Price        float64       `json:"price" gorm:"not null"`
	Type         EstateType    `json:"type" gorm:"not null" validate:"required,estate_type"`
	InRent       bool          `json:"in_rent" gorm:"default:false"`
	User         User          `json:"-" gorm:"foreignKey:UserID" validate:"-"`
	UserID       uint          `json:"user_id" validate:"gte=0"`
	Transactions []Transaction `gorm:"foreignKey:EstateID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Custom validation function for EstateType
func ValidateEstateType(fl validator.FieldLevel) bool {
	_, ok := fl.Field().Interface().(EstateType)
	if !ok {
		return false
	}
	return true

	// Check if the value exists in the slice, using only part of EstateType
	// for _, validType := range ValidEstateTypes {
	//     if estateType == validType {
	//         return true
	//     }
	// }
	// return false
}

func EstateValidate(estate *Estate) error {
	validate := validator.New()
	validate.RegisterValidation("estate_type", ValidateEstateType)
	return validate.Struct(estate)
}
