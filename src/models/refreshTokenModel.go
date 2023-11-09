package models

import (
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model

	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"not null"`
	ExpiresAt int64  `gorm:"not null"`
}
