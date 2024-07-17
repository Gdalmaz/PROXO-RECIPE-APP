package controllers

import (
	"errors"
	"proxo-go-application/config"
	"proxo-go-application/database"
	"proxo-go-application/helpers"
	"proxo-go-application/models"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddTaste(c *fiber.Ctx) error {
	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}

	food := new(models.Food)
	err := c.BodyParser(&food)

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "not loading file", "data": err.Error()})
	}

	fileBytes, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "not reading file", "data": err.Error()})
	}

	defer fileBytes.Close()

	imageBytes := make([]byte, file.Size)
	_, err = fileBytes.Read(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to read file", "data": err.Error()})
	}

	id, url, err := config.CloudConnect(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to upload cloud", "data": err.Error()})
	}
	food.Image = id
	food.ImageUrl = url

	food.UserID = sessionUser.ID
	foodname := c.FormValue("foodname")
	materials := c.FormValue("materials")
	eatpersonStr := c.FormValue("eatperson")
	specification := c.FormValue("specification")
	guesspriceStr := c.FormValue("guessprice")
	preparationtimeStr := c.FormValue("preparationtime")

	if len(foodname) != 0 {
		food.FoodName = foodname
	}

	if len(materials) != 0 {
		food.Materials = materials
	}

	if eatpersonStr != "" {
		eatperson, err := strconv.Atoi(eatpersonStr)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Invalid bus price", "data": err.Error()})
		}
		if eatperson > 0 {
			food.EatPerson = eatperson
		} else {
			return c.Status(402).JSON(fiber.Map{"status": "error", "message": "Bus price must be greater than zero"})
		}
	}

	if len(specification) != 0 {
		food.Specification = specification
	}

	if guesspriceStr != "" {
		guessprice, err := strconv.Atoi(guesspriceStr)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Invalid bus price", "data": err.Error()})
		}
		if guessprice > 0 {
			food.GuessPrice = guessprice
		} else {
			return c.Status(402).JSON(fiber.Map{"status": "error", "message": "Bus price must be greater than zero"})
		}
	}
	if preparationtimeStr != "" {
		preparationtime, err := strconv.Atoi(preparationtimeStr)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Invalid bus price", "data": err.Error()})
		}
		if preparationtime > 0 {
			food.PreparationTime = preparationtime
		} else {
			return c.Status(402).JSON(fiber.Map{"status": "error", "message": "Bus price must be greater than zero"})
		}
	}
	err = database.DB.Db.Create(&food).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error create food step", "data": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your food created successfully", "data": sessionUser})
}

func UpdateTaste(c *fiber.Ctx) error {
	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}
	id := c.Params("id")
	var foods models.Food
	err := database.DB.Db.First(&foods, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "foods not found"})
	}
	if sessionUser.ID != foods.UserID {
		return c.Status(402).JSON(fiber.Map{"status": "error", "message": "you don't have permission for this"})
	}
	updateData := foods
	err = c.BodyParser(&updateData)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error bodyparsing step"})
	}
	err = database.DB.Db.Model(foods).Updates(updateData).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error update step"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your foods updated successfully", "data": foods})

}

func DeleteTaste(c *fiber.Ctx) error {
	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}
	id := c.Params("id")
	UserID := sessionUser.ID
	food := new(models.Food)
	err := database.DB.Db.Where("id=?", id).First(&food).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error not found the food post"})
	}
	if food.UserID != UserID {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "permission denied"})
	}
	click := new(models.Popularity)

	err = database.DB.Db.Where("food_id", id).Delete(&click).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error delete click step"})
	}
	err = database.DB.Db.Delete(&food).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error delete step", "data": err.Error()})
	}
	err = config.DeleteClickCountFromRedis(uint(food.ID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error deleting click count from Redis"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your food post deleted successfully"})
}

func GetAllTaste(c *fiber.Ctx) error {
	_, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}

	var foods []models.Food
	if err := database.DB.Db.Preload("User").Find(&foods).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "failed to fetch food items"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "successfully fetched foods", "data": foods})
}

func GetAllYourTaste(c *fiber.Ctx) error {
	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}
	var foods []models.Food
	err := database.DB.Db.Preload("User").Where("user_id=?", sessionUser.ID).Order("id DESC").Find(&foods).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error loading screen", "data": err})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "loading successfull", "data": foods})

}

func GetClickTaste(c *fiber.Ctx) error {
	sessionUser, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}

	var foods models.Food
	err := c.BodyParser(&foods)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error foods bodyparser", "data": err.Error()})
	}

	err = database.DB.Db.Preload("User").Where("id=?", foods.ID).Find(&foods).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "not found foods", "data": err.Error()})
	}

	var click models.Popularity
	err = c.BodyParser(&click)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error click bodyparser", "data": err.Error()})
	}

	err = database.DB.Db.Where("food_id = ? AND user_id = ?", foods.ID, sessionUser.ID).First(&click).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			click.UserID = sessionUser.ID
			click.FoodID = foods.ID
			click.ClickNumber = 1 //ilk tıklama
		} else {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error finding click record", "data": err.Error()})
		}
	} else {
		click.ClickNumber++
		if click.ClickNumber > 5 {
			err = config.SaveClickCountToRedis(uint(foods.ID), click.ClickNumber)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error saving to Redis", "data": err.Error()})
			}
		}
	}
	err = database.DB.Db.Save(&click).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error save click step", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": foods})
}
func PopularTaste(c *fiber.Ctx) error {
	_, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}
	clickcount, err := config.GetAllClickCountsFromRedis()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error click count"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": clickcount})
}

//Filtreleme özelliği
func SearchHandler(c *fiber.Ctx) error{
	_, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "user not found in context"})
	}
	db := database.DB.Db
	searchText := c.Query("q")
	items ,err := helpers.SearchItems(db ,searchText)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status":"error","message":"error filter","err":err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"status":"success","message":"success fiter","data":items})
}