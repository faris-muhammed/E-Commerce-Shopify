package model

import (
	"time"

	"gorm.io/gorm"
)

type AdminModel struct {
	Id       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}
type Category struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDeleted   bool   `gorm:"default:false"`
}
type OTPDetails struct {
	Id        uint
	Email     string
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}
