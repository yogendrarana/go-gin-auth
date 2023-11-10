package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	RefreshTokens []RefreshToken
}

// Get user by ID
func GetUserByID(db *gorm.DB, userID uint) (*User, error) {
	var user User
	err := db.Preload("RefreshTokens").Where("id = ?", userID).First(&user).Error
	return &user, err
}
