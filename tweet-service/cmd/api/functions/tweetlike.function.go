package functions

import (
	"context"
	"strconv"
	"time"
	"tweet/cmd/api/models"
	"tweet/cmd/cache"
	"tweet/cmd/database"
	rabbitmq "tweet/cmd/rabbitMQ"

	"github.com/goccy/go-json"
)

var ctx = context.Background()

func LikeTweet(userId int64, tweetId int64) (*models.FavouriteTweet, error) {
	db := database.DB

	tweetLike := models.FavouriteTweet{
		TweetId:      int64(tweetId),
		UserId:       int64(userId),
		CreationTime: time.Now(),
	}

	err := db.Create(&tweetLike).Error
	if err != nil {
		return nil, err
	}

	var publishPayload = rabbitmq.LikePayload{
		TweetId: tweetId,
		UserId:  userId,
		IsLike:  true,
	}

	jsonMarshal, err := json.Marshal(publishPayload)

	if err != nil {
		return nil, err
	}

	err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME, "LIKE", string(jsonMarshal))

	if err != nil {
		return nil, err
	}

	return &tweetLike, nil
}

func UnlikeTweet(userId int64, tweetId int64) (*models.FavouriteTweet, error) {
	db := database.DB

	tweetLike := models.FavouriteTweet{
		TweetId:      int64(tweetId),
		UserId:       int64(userId),
		CreationTime: time.Now(),
	}

	err := db.Where("tweet_id = ? AND user_id = ?", tweetId, userId).Delete(&models.FavouriteTweet{}).Error
	if err != nil {
		return nil, err
	}
	var publishPayload = rabbitmq.LikePayload{
		TweetId: tweetId,
		UserId:  userId,
		IsLike:  false,
	}

	jsonMarshal, err := json.Marshal(publishPayload)

	if err != nil {
		return nil, err
	}

	err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME, "LIKE", string(jsonMarshal))

	if err != nil {
		return nil, err
	}

	return &tweetLike, nil
}

func UserLikeTweets(userId int64) ([]string, error) {
	db := database.DB

	var userLiked []string

	userCacheKey := "userLike:" + strconv.FormatInt(userId, 10)

	cacheGet := cache.RedisClient.HKeys(ctx, userCacheKey)

	if cacheGet.Err() == nil {
		return cacheGet.Val(), nil
	}

	err := db.Model(&models.FavouriteTweet{}).Where("user_id = ?", userId).Select("tweet_id").Find(&userLiked).Error
	if err != nil {
		return nil, err
	}

	return userLiked, nil
}
