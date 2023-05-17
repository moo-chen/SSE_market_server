package model

import (
	"time"
)

// Sue [...]
type Sue struct {
	Sueid      int       `gorm:"autoIncrement:true;primaryKey;column:sueID;type:int;not null" json:"-"`
	Targettype string    `gorm:"column:targettype;type:enum('post','pcomment','ccomment');not null" json:"targettype"`
	Targetid   int       `gorm:"column:targetID;type:int;not null" json:"targetId"`
	Userid     int       `gorm:"index:psueuser;column:userID;type:int;not null" json:"userId"`
	User       User      `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
	Reason     string    `gorm:"column:reason;type:varchar(1000);default:null" json:"reason"`
	Time       time.Time `gorm:"column:time;type:datetime;default:null" json:"time"`
	Status     string    `gorm:"column:status;type:enum(' ok','nosin','wait');default:null" json:"status"`
	Finish     bool      `gorm:"column:finish;type:tinyint(1);default:null" json:"finish"`
}
