package model

import "time"

type User struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	Username  string   `json:"username" gorm:"unique;not null" validate:"required"`
	Names     string   `json:"names" validate:"required"`
	Email     string   `json:"email" gorm:"unique;not null" validate:"required"`
	Phone     string   `json:"phone" gorm:"unique;not null" validate:"required"`
	Password  string   `json:"password" gorm:"not null" validate:"required"`
	Estates   []Estate `json:"estates" gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
