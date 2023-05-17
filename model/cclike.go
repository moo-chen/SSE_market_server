package model

// Cclike [...]
type Cclike struct {
	Cclikeid   int      `gorm:"primaryKey;column:cclikeID;type:int;not null" json:"-"`
	Cctargetid int      `gorm:"index:ccliketarget;column:cctargetID;type:int;not null" json:"cctargetId"`
	Ccomment   Ccomment `gorm:"joinForeignKey:cctargetID;foreignKey:ccommentID;references:Cctargetid" json:"ccommentList"`
	Userid     int      `gorm:"index:cclikeuser;column:userID;type:int;not null" json:"userId"`
	User       User     `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
}
