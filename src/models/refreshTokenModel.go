package models

import (
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model

	UserID    uint   `gorm:"not null"`
	TokenHash string `gorm:"not null"`
	ExpiresAt int64  `gorm:"not null"`
}

// Get refresh token by user ID
func GetRefreshTokenByUserID(db *gorm.DB, userID uint) (*RefreshToken, error) {
	var refreshToken RefreshToken
	err := db.Where("user_id = ?", userID).First(&refreshToken).Error
	return &refreshToken, err
}
