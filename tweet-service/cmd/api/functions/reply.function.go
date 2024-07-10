package functions

import (
	"time"
	"tweet/cmd/api/models"
	"tweet/cmd/database"

	"gorm.io/gorm"
)

func CreateReply(reply models.Reply) (*models.Reply, error) {
	db := database.DB

	replyBody := models.Reply{
		UserId:       int(reply.UserId),
		CreationTime: time.Now(),
		Content:      reply.Content,
		TotalLikes:   0,
		MediaFiles:   reply.MediaFiles,
		TweetId:      int(reply.TweetId),
	}

	if err := db.Model(models.Reply{}).Create(&replyBody).Error; err != nil {
		return nil, err
	}

	return &replyBody, nil
}

func GetAllReplies(tweetId int64, page int64, size int64) (*[]models.Reply, error) {
	db := database.DB

	var replies *[]models.Reply

	offset := (page - 1) * size

	if err := db.Model(models.Reply{}).Preload("User").Where("tweet_id = ?", tweetId).Limit(int(size)).Offset(int(offset)).Order("creation_time DESC").Find(&replies).Error; err != nil {
		return nil, err
	}

	return replies, nil
}

func LikeReplyAction(replyId int64, userId int64, isLike bool) error {
	db := database.DB

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if isLike {
		if err := LikeReplyUser(replyId, userId); err != nil {
			return err
		}

		if err := db.Model(&models.Reply{}).Where("total_likes >= 0 AND id = ?", replyId).Update("total_likes", gorm.Expr("total_likes + ?", 1)).Error; err != nil {
			return err
		}
	} else {      
		if err := DisLikeReplyUser(replyId, userId); err != nil {
			return err
		}
		if err := db.Model(&models.Reply{}).Where("total_likes > 0 AND id = ?", replyId).Update("total_likes", gorm.Expr("total_likes - ?", 1)).Error; err != nil {
			return err
		}
	}

	return nil
}
