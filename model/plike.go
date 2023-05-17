package model

// Plike [...]
type Plike struct {
	Plikeid   int  `gorm:"primaryKey;column:plikeID;type:int;not null" json:"-"`
	Ptargetid int  `gorm:"index:pliketarget;column:ptargetID;type:int;not null" json:"ptargetId"`
	Post      Post `gorm:"joinForeignKey:ptargetID;foreignKey:postID;references:Ptargetid" json:"postList"`
	Userid    int  `gorm:"index:plikeuser;column:userID;type:int;not null" json:"userId"`
	User      User `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
}
