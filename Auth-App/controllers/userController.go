package controllers

import (
	"auth/database"
	"auth/helpers"
	"auth/middleware"
	"auth/models"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error to bodyparsing step"})
	}
	controlMail, _ := helpers.MailControl(user.Mail)
	if controlMail == true {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "you already have this account on this mail"})
	}

	controlPhone, _ := helpers.PhoneControl(user.Phone)
	if controlPhone == true {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "you already have this account on this phone number"})
	}
	user.Password = helpers.HashPass(user.Password)
	err = database.DB.Db.Create(user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error create account step", "data": err})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your account creating successfully", "data": user})
}

func LogIn(c *fiber.Ctx) error {
	user := new(models.User)
	successUser := new(models.LogIn)
	err := c.BodyParser(&successUser)
	successUser.LogPassword = helpers.HashPass(successUser.LogPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error to bodyparsing step"})
	}
	err = database.DB.Db.Where("mail =? and password=?", successUser.LogMail, successUser.LogPassword).First(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "wrong password or mail", "data": err})
	}

	token, err := middleware.CreateToken(user.FirstName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "not created token"})
	}

	session := new(models.Session)
	session.UserID = user.ID
	session.Token = token
	database.DB.Db.Create(&session)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "login successfully", "data": user})
}

func UpdatePassword(c *fiber.Ctx) error {

	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Fiber bağlamında kullanıcı bulunamadı",
		})
	}

	// Güncellenecek şifreyi alın
	updatePassword := new(models.UpdatePassword)
	if err := c.BodyParser(updatePassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Şifre güncelleme isteği işlenirken bir hata oluştu",
		})
	}

	// Eski şifreyi kontrol edin
	if helpers.HashPass(updatePassword.OldPassword) != sessionUser.Password {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Eski şifreyi yanlış girdiniz",
		})
	}

	// Yeni şifrelerin eşleştiğini kontrol edin
	if updatePassword.NewPassword1 != updatePassword.NewPassword2 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Yeni şifreler eşleşmiyor",
		})
	}

	// Yeni şifreyi hashleyin ve kullanıcı nesnesine kaydedin
	sessionUser.Password = helpers.HashPass(updatePassword.NewPassword1)
	if err := database.DB.Db.Save(&sessionUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Şifre güncellenirken bir hata oluştu",
			"data":    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Şifreniz başarıyla güncellendi",
	})
}
func DeleteAccount(c *fiber.Ctx) error {
	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}

	err := database.DB.Db.Exec("DELETE FROM sessions WHERE user_id = ?", sessionUser.ID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error delete session step1111111111111"})
	}
	err = database.DB.Db.Where("id=?", sessionUser.ID).Delete(&sessionUser).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error delete account step"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your account deleted successfully"})
}

func LogOut(c *fiber.Ctx) error {
	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}
	err := database.DB.Db.Exec("DELETE FROM sessions WHERE user_id = ?", sessionUser.ID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error delete session step"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "logout successfully"})
}
