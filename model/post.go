package model

import (
	"github.com/jinzhu/gorm"
 	"time"
)

type Post struct {
    gorm.Model
    UserID    int  `gorm:"not null"`
    Partition string `gorm:"not null"`
    Title     string `gorm:"type:varchar(15);not null"`
    Content   string `gorm:"type:varchar(5000);not null"`
    Like      int    `gorm:"not null"`
    Comment   int    `gorm:"not null"`
    Heat      float64 `gorm:"not null"`
    PostTime  time.Time `gorm:"not null"`
}
