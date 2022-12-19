package entity

import (
	"time"
)

type Status int

const (
	UserStatusPendingActive = Status(0)
	UserStatusActive        = Status(1)
)

type User struct {
	UserID   string `gorm:"primaryKey"`
	Mail     string `gorm:"index;unique"`
	Password string
}

type RoleID int

const (
	RoleManageBook = RoleID(1)
	RoleManageUser = RoleID(2)
)

type UserRole struct {
	ID     uint   `gorm:"primaryKey"`
	UserID string `gorm:"index"`
	RoleID RoleID
}

type UserStatus struct {
	UserID string `gorm:"primaryKey"`
	Status Status
}

type UserActive struct {
	Code      string `gorm:"primaryKey"`
	CreatedAt time.Time
	UserID    string
	ExpiredAt time.Time
}

type UserResetPassword struct {
	Code      string `gorm:"primaryKey"`
	CreatedAt time.Time
	UserID    string
	ExpiredAt time.Time
}
