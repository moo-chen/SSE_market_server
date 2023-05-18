package model

// Plike [...]
type Plike struct {
	PlikeID   int  `gorm:"primary_key;column:plikeID"`
	PtargetID int  `gorm:"index:pliketarget;column:ptargetID;type:int;not null"`
	Post      Post `gorm:"association_foreignkey:ptargetID;foreignkey:postID"`
	UserID    int  `gorm:"index:plikeuser;column:userID;type:int;not null"`
	User      User `gorm:"association_foreignkey:userID;foreignkey:userID"`
}
