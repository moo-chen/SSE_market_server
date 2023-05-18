package model

// Notice [...]
type Notice struct {
	NoticeID int    `gorm:"primary_key;column:noticeID"`
	Receiver int    `gorm:"index:noticereceiver;column:receiver;type:int;not null"`
	User     User   `gorm:"association_foreignkey:receiver;foreignkey:userID"`
	Sender   int    `gorm:"index:noticesender;column:sender;type:int;not null"`
	Type     string `gorm:"column:type;type:enum('pcomment','ccomment','punish')"`
	Ntext    string `gorm:"column:ntext;type:varchar(100)"`
	Read     bool   `gorm:"column:read;type:tinyint(1)"`
}
