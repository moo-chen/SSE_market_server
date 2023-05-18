package model

import (
	"time"
)

// Feedback [...]
type Feedback struct {
	FeedbackID int       `gorm:"primary_key;column:feedbackID"`
	UserID     int       `gorm:"index:feedbackuser;column:userID;type:int"`
	User       User      `gorm:"association_foreignkey:userID;foreignkey:userID"`
	Ftext      string    `gorm:"column:ftext;type:varchar(1000)"`
	Time       time.Time `gorm:"column:time;type:datetime"`
	Status     string    `gorm:"column:status;type:enum('ok','wait')"`
}
