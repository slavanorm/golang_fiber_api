package model

import (
	"time"
)

type Transaction struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"-" gorm:"foreignKey:UserID" validate:"-"`
	EstateID   uint      `json:"estate_id" validate:"required"`
	Amount     float64   `json:"amount" gorm:"not null" validate:"required"`
	Repetitive bool      `json:"repetitive" gorm:"default:false"`
	DateActual time.Time `json:"date_actual" gorm:"not null"`
	DatePeriod time.Time `json:"date_period,omitempty"`
	Comment1   string    `json:"comment1,omitempty"`
	Comment2   string    `json:"comment2,omitempty"`
	Comment3   string    `json:"comment3,omitempty"`
	CreatedAt time.Time 
    UpdatedAt time.Time 
}
