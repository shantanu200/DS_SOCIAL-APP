package cachefunction

import (
	"strconv"
	"timeline/cmd/api/models"
	"timeline/cmd/cache"
)

func GetAllUserTweets(userId int64) ([]string, error) {
	cacheKey := "user:tweets:" + strconv.FormatInt(userId, 10)

	tweets, err := cache.RedisClient.ZRevRange(ctx, cacheKey, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func GetUserLikedTweets(userId int64) ([]models.Tweet, error) {

	cacheKey := "userLike:" + strconv.FormatInt(userId, 10)

	tweetIds, err := cache.RedisClient.SMembers(ctx, cacheKey).Result()

	if err != nil {
		return nil, err
	}

	tweets, err := GetTweets(tweetIds)

	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func GetUserPostTweets(userId int64) ([]models.Tweet, error) {
	cacheKey := "user:tweets:" + strconv.FormatInt(userId, 10)

	tweetIds, err := cache.RedisClient.ZRevRange(ctx, cacheKey, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	tweets, err := GetTweets(tweetIds)

	if err != nil {
		return nil, err
	}

	return tweets, nil
}
