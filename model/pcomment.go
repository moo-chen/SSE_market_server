package model

import (
	"time"
)

// Pcomment [...]
type Pcomment struct {
	Pcommentid int       `gorm:"primaryKey;column:pcommentID;type:int;not null" json:"-"`
	Userid     int       `gorm:"index:pcommentuser;column:userID;type:int;not null" json:"userId"`
	User       User      `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
	Ptargetid  int       `gorm:"index:pcommenttarget;column:ptargetID;type:int;not null" json:"ptargetId"`
	Post       Post      `gorm:"joinForeignKey:ptargetID;foreignKey:postID;references:Ptargetid" json:"postList"`
	LikeNum    int       `gorm:"column:like_num;type:int;default:null" json:"likeNum"`
	Pctext     string    `gorm:"column:pctext;type:varchar(1000);default:null" json:"pctext"`
	Time       time.Time `gorm:"column:time;type:datetime;default:null" json:"time"`
}
