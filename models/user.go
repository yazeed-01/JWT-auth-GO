package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:'user'"`
	IsOnline bool   `gorm:"default:false"`
}

type AccessToken struct {
	gorm.Model
	UserID uint   `gorm:"not null"`
	Token  string `gorm:"not null"`
	IsDead bool   `gorm:"default:false"`
}

type RefreshToken struct {
	gorm.Model
	UserID uint   `gorm:"not null"`
	Token  string `gorm:"not null"`
	IsDead bool   `gorm:"default:false"`
}
