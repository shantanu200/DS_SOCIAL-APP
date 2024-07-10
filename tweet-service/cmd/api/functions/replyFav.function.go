package functions

import (
	"strconv"
	"time"
	"tweet/cmd/api/models"
	"tweet/cmd/cache"
	"tweet/cmd/database"
	rabbitmq "tweet/cmd/rabbitMQ"

	"github.com/goccy/go-json"
)

func LikeReplyUser(replyId int64, userId int64) error {
	db := database.DB

	replyBody := models.FavouriteReply{
		ReplyId:      int(replyId),
		UserId:       int(userId),
		CreationTime: time.Now(),
	}

	if err := db.Model(models.FavouriteReply{}).Create(&replyBody).Error; err != nil {
		return err
	}

	var rabbitError error

	go func() {

		publishPayload := rabbitmq.ReplyLikePayload{
			ReplyId: replyId,
			UserId:  userId,
			IsLike:  true,
		}

		jsonMarshal, err := json.Marshal(publishPayload)

		if err != nil {
			rabbitError = err
		}

		err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME, "LIKE_REPLY", string(jsonMarshal))

		if err != nil {
			rabbitError = err
		}
	}()

	if rabbitError != nil {
		return rabbitError
	}

	return nil
}

func DisLikeReplyUser(replyId int64, userId int64) error {
	db := database.DB

	if err := db.Where("reply_id = ? AND user_id = ?", replyId, userId).Delete(&models.FavouriteReply{}).Error; err != nil {
		return err
	}

	var rabbitError error

	go func() {
		publishPayload := rabbitmq.ReplyLikePayload{
			ReplyId: replyId,
			UserId:  userId,
			IsLike:  false,
		}

		jsonMarshal, err := json.Marshal(publishPayload)

		if err != nil {
			rabbitError = err
		}

		err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME, "LIKE_REPLY", string(jsonMarshal))

		if err != nil {
			rabbitError = err
		}
	}()

	if rabbitError != nil {
		return rabbitError
	}

	return nil
}

func UserLikeRepliesIds(userId int64) (*[]string, error) {
	userReplyLikeKey := "userReplyLike:" + strconv.FormatInt(int64(userId), 10)

	replyIds, err := cache.RedisClient.ZRange(ctx, userReplyLikeKey, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	return &replyIds, nil
}
