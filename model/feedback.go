package model

import (
	"time"
)

// Feedback [...]
type Feedback struct {
	Feedbackid int       `gorm:"primaryKey;column:feedbackID;type:int;not null" json:"-"`
	Userid     int       `gorm:"index:feedbackuser;column:userID;type:int;default:null" json:"userId"`
	User       User      `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
	Ftext      string    `gorm:"column:ftext;type:varchar(1000);default:null" json:"ftext"`
	Time       time.Time `gorm:"column:time;type:datetime;default:null" json:"time"`
	Status     string    `gorm:"column:status;type:enum('ok','wait');default:null" json:"status"`
}
