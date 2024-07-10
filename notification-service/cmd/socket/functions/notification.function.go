package functions

import (
	"fmt"
	"notification/cmd/database"
	"notification/cmd/socket/models"
	"strings"
	"time"
)

const BATCH_SIZE = 100

func CreateNotifications(doc []models.Notification) error {
	db := database.DB

	for idx := range doc {
		if doc[idx].Type != "LIKE" {
			doc[idx].PostId = 0
		}
		doc[idx].CreationTime = time.Now().UTC()
		doc[idx].Read = false
	}

	err := db.CreateInBatches(doc, BATCH_SIZE).Error

	if err != nil {
		if exists := strings.Contains(err.Error(), "duplicate key"); exists {
			fmt.Println(err)
		} else {
			return err
		}
	}

	return nil
}

func ReadNotifications(userId int64, notifications []int64) error {
	db := database.DB

	err := db.Model(models.Notification{}).Where("id IN ?", notifications).Updates(map[string]interface{}{"read": true}).Error

	if err != nil {
		return err
	}

	return nil
}

func GetAllUserNotifications(userId int64) (*[]models.Notification, error) {
	db := database.DB

	var notifications *[]models.Notification

	err := db.Model(models.Notification{}).Preload("Author").Where("user_id = ?", userId).Order("read").Find(&notifications).Error

	if err != nil {
		return nil, err
	}

	return notifications, nil
}
