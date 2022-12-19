package entity

import (
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
	Token     string `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	UserID    string `gorm:"index"`
	ExpiredAt time.Time
}
