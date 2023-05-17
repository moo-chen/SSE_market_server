package model

import (
	"time"
)

// Ccomment [...]
type Ccomment struct {
	Ccommentid int       `gorm:"primaryKey;column:ccommentID;type:int;not null" json:"-"`
	Userid     int       `gorm:"index:ccommentuser;column:userID;type:int;not null" json:"userId"`
	User       User      `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
	Ctargetid  int       `gorm:"index:ccommenttarget;column:ctargetID;type:int;not null" json:"ctargetId"`
	Pcomment   Pcomment  `gorm:"joinForeignKey:ctargetID;foreignKey:pcommentID;references:Ctargetid" json:"pcommentList"`
	LikeNum    int       `gorm:"column:like_num;type:int;default:null" json:"likeNum"`
	Cctext     string    `gorm:"column:cctext;type:varchar(100);default:null" json:"cctext"`
	Time       time.Time `gorm:"column:time;type:datetime;default:null" json:"time"`
}
