package model

// Admin [...]
type Admin struct {
	Account  string `gorm:"primaryKey;column:account;type:varchar(100);not null" json:"-"`
	Password string `gorm:"column:password;type:varchar(20);default:null" json:"password"`
}
