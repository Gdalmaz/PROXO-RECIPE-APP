package models

type User struct {
	ID        int    `json"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type LogIn struct {
	LogMail     string `json:"logmail"`
	LogPassword string `json:"logpassword"`
}

type UpdatePassword struct {
	OldPassword  string `json:"oldpassword"`
	NewPassword1 string `json:"newpassword1"`
	NewPassword2 string `json:"newpassword2"`
}
