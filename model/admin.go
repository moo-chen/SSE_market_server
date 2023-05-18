package model

// Admin [...]
type Admin struct {
	Account  string `gorm:"primary_key;column:account;type:varchar(100);not null"`
	Password string `gorm:"column:password;type:varchar(20)"`
}
