package models

type SendMailInApp struct {
	Mail     string `json:"mail"`
	SendMail string `json:"sendmail"`
	Text     string `json:"text"`
}
