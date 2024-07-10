package functions

import (
	"time"
	"tweet/cmd/api/models"
	"tweet/cmd/database"
)

func CreateThread(thread *models.Thread) (*models.Thread, error) {
	db := database.DB

	threadBody := models.Thread{
		UserId:       thread.UserId,
		ReplyId:      thread.ReplyId,
		ThreadId:     thread.ReplyId,
		TweetId:      thread.TweetId,
		Content:      thread.Content,
		MediaFiles:   thread.MediaFiles,
		CreationTime: time.Now(),
		TotalLikes:   0,
	}

	if err := db.Model(models.Thread{}).Create(&threadBody).Error; err != nil {
		return nil, err
	}

	return &threadBody, nil
}

func GetAllThreadReplies(threadId int64, page int64, size int64) (*[]models.Thread, error) {
	db := database.DB
    
	var threads *[]models.Thread

	offset := (page - 1) * size

	if err := db.Model(models.Thread{}).Preload("User").Where("thread_id = ?",threadId).Limit(int(size)).Offset(int(offset)).Order("creation_time DESC").Find(&threads).Error; err != nil {
		return nil,err
	}

	return threads,nil
}

