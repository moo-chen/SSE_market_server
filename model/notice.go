package model

// Notice [...]
type Notice struct {
	Noticeid int    `gorm:"primaryKey;column:noticeID;type:int;not null" json:"-"`
	Receiver int    `gorm:"index:noticereceiver;column:receiver;type:int;not null" json:"receiver"`
	User     User   `gorm:"joinForeignKey:receiver;foreignKey:userID;references:Receiver" json:"userList"`
	Sender   int    `gorm:"index:noticesender;column:sender;type:int;not null" json:"sender"`
	Type     string `gorm:"column:type;type:enum('pcomment','ccomment','punish');default:null" json:"type"`
	Ntext    string `gorm:"column:ntext;type:varchar(100);default:null" json:"ntext"`
	Read     bool   `gorm:"column:read;type:tinyint(1);default:null" json:"read"`
}
