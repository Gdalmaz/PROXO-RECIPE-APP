package main

import (

	"log"
	"mail-services/config"
	"mail-services/database"
	"mail-services/models"

)


func main() {
	var sendmail models.SendMailInApp
	sender := sendmail.SendMail
	text := sendmail.Text
	err := config.RabbitMqPublish([]byte(text), sender)
	if err != nil {
		log.Println("Not publish")
		panic(err)
	}
	err = config.RabbitMqConsume(sender, sendmail.Mail)
	if err != nil {
		log.Println("Not consume")
		panic(err)
	}
	err = database.DB.Db.Save(&sendmail).Error
	if err != nil {
		log.Println("Kayıt yapılamadı")
		panic(err)
	}
	log.Println("Mail Gönderildi Ve Log Kaydı yapıldı")



}
