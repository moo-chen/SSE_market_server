package model

import "github.com/jinzhu/gorm"

type Like struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	PostID uint `gorm:"not null"`
}