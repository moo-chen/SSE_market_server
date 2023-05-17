package model

// User [...]
type User struct {
	Userid    int    `gorm:"primaryKey;column:userID;type:int;not null" json:"-"`
	Phone     string `gorm:"column:phone;type:char(20);default:null" json:"phone"`
	Email     string `gorm:"column:email;type:varchar(255);default:null" json:"email"`
	Password  string `gorm:"column:password;type:varchar(255);default:null" json:"password"`
	Name      string `gorm:"column:name;type:varchar(50);default:null" json:"name"`
	Num       int    `gorm:"column:num;type:int;default:null" json:"num"`
	Profile   string `gorm:"column:profile;type:varchar(100);default:null" json:"profile"`
	Intro     string `gorm:"column:intro;type:varchar(255);default:null" json:"intro"`
	IDpass    bool   `gorm:"column:idpass;type:tinyint(1);default:null" json:"idpass"`
	Ban       bool   `gorm:"column:ban;type:tinyint(1);default:null" json:"ban"`
	Punishnum int    `gorm:"column:punishnum;type:int;default:null" json:"punishnum"`
}
