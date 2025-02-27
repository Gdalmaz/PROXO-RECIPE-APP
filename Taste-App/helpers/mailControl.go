package helpers

import (
	"proxo-go-application/database"
	"proxo-go-application/models"
)

func MailControl(mail string) (bool, error) {
	user := new(models.User)

	err := database.DB.Db.Where("mail=?", mail).First(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
