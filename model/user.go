package model

type User struct {
	ID        uint   
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}
