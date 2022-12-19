package entity

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	ID     uint
	Name   string `gorm:"index;unique"`
	Author string `gorm:"index"`
}

type BookRent struct {
	gorm.Model        //用gorm的软删除特性来进行归还的确认
	BookID     uint   `gorm:"index"`
	RentBy     string `gorm:"index"`
	ExpiredAt  time.Time
}
