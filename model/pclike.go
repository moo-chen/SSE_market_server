package model

// Pclike [...]
type Pclike struct {
	Pclikeid   int      `gorm:"primaryKey;column:pclikeID;type:int;not null" json:"-"`
	Pctargetid int      `gorm:"index:pcliketarget;column:pctargetID;type:int;not null" json:"pctargetId"`
	Pcomment   Pcomment `gorm:"joinForeignKey:pctargetID;foreignKey:pcommentID;references:Pctargetid" json:"pcommentList"`
	Userid     int      `gorm:"index:pclikeuser;column:userID;type:int;not null" json:"userId"`
	User       User     `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
}
