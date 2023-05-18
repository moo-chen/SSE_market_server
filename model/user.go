package model

// User [...]
type User struct {
	UserID    int    `gorm:"primary_key;column:userID"`
	Phone     string `gorm:"column:phone;type:char(20)"`
	Email     string `gorm:"column:email;type:varchar(255)"`
	Password  string `gorm:"column:password;type:varchar(255)"`
	Name      string `gorm:"column:name;type:varchar(50)"`
	Num       int    `gorm:"column:num;type:int"`
	Profile   string `gorm:"column:profile;type:varchar(100)"`
	Intro     string `gorm:"column:intro;type:varchar(255)"`
	IDpass    bool   `gorm:"column:idpass;type:tinyint(1)"`
	Ban       bool   `gorm:"column:ban;type:tinyint(1)"`
	Punishnum int    `gorm:"column:punishnum;type:int"`
}
