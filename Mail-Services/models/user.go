package models

type User struct {
	ID        int    `json"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}
