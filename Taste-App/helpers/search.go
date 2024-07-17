package helpers

import (
	"proxo-go-application/database"
	"proxo-go-application/models"

	"gorm.io/gorm"
)

func SearchItems(db *gorm.DB, searchText string) ([]models.Food, error) {

	var food []models.Food

	searchText = "%" + searchText + "%"
	err := database.DB.Db.Where("food_name ILIKE ?", searchText).Find(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil

}