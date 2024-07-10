package dbhandler

import (
	"timeline/cmd/api/models"
	"timeline/cmd/database"
)

func GetAllTweetsForUser(tweetIds []string) ([]models.Tweet, error) {
	db := database.DB
	var tweets []models.Tweet

	err := db.Model(&models.Tweet{}).Preload("User").Where("id IN (?) ", tweetIds).Find(&tweets).Error

	if err != nil {
		return nil, err
	}

	return tweets, nil
}
