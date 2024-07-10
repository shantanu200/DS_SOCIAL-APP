package cachefunction

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	"timeline/cmd/api/models"
	"timeline/cmd/cache"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func GetUserTimelineCache(userId int64, page int64, size int64) ([]string, error) {
	cacheKey := "timeline:user:" + strconv.FormatInt(userId, 10)
	var start int64
	if page == 1 {
		start = 0
	} else {
		start = (page-1)*size + 1
	}

	end := (page)*size - 1

	log.Printf("[TimeLine Fetch] %d page and %d size ", start, end)

	tweetIds, err := cache.RedisClient.ZRevRange(ctx, cacheKey, start, end).Result()

	log.Println(tweetIds)

	if err != nil {
		return nil, err
	}

	return tweetIds, nil
}

func GetHomeTimelineCache(userId int64, page int64, size int64) ([]models.Tweet, error) {
	cacheKey := "tweets"

	tweets, err := cache.RedisClient.HRandFieldWithValues(ctx, cacheKey, 10).Result()

	var Tweet models.Tweet
	var Tweets []models.Tweet
	for _, tweet := range tweets {
		err := json.Unmarshal([]byte(tweet.Value), &Tweet)

		if err != nil {
			log.Fatalln(err)
		}

		Tweets = append(Tweets, Tweet)
	}

	if err != nil {
		return nil, err
	}

	return Tweets, nil
}

func GetTweets(tweetIds []string) ([]models.Tweet, error) {
	cacheKey := "tweets"
	var results []models.Tweet
	var tweet *models.Tweet

	for _, tweetId := range tweetIds {
		cacheTweet, err := cache.RedisClient.HGet(ctx, cacheKey, tweetId).Result()
		if err != nil {
			fmt.Printf("Unable to fetch tweet with %s ", tweetId)
		}

		err = json.Unmarshal([]byte(cacheTweet), &tweet)
		if err != nil {
			fmt.Printf("Unable to unmarshal tweet with %s ", tweetId)
		}

		results = append(results, *tweet)
	}

	return results, nil
}

func UpdateTweets(tweets []models.Tweet) error {
	cacheKey := "tweets"

	for _, tweet := range tweets {
		jsonMarshal, err := json.Marshal(tweet)

		if err != nil {
			log.Fatalf("[UPDATE TWEET ERROR] %s", err.Error())
		}

		err = cache.RedisClient.HSet(ctx, cacheKey, tweet.Id, string(jsonMarshal)).Err()

		if err != nil {
			log.Fatalf("[UPDATE TWEET ERROR] %s", err.Error())
		}
	}

	return nil
}

func RemoveTweetsFromTimeline(removeUserId int64, userId int64) error {
	userTweetKey := "user:tweets:" + strconv.FormatInt(removeUserId, 10)

	tweetsId, err := cache.RedisClient.SMembers(ctx, userTweetKey).Result()

	if err != nil {
		return err
	}

	cacheKey := "timeline:user:" + strconv.FormatInt(userId, 10)

	err = cache.RedisClient.ZRem(ctx, cacheKey, tweetsId).Err()

	if err != nil {
		return err
	}

	return nil
}

func UpdateUserTimeLineForFollowedUser(followerId int64, userId int64) error {
	userTweetKey := "user:tweets:" + strconv.FormatInt(followerId, 10)

	tweetsId, err := cache.RedisClient.ZRevRange(ctx, userTweetKey, 0, 3).Result()

	if err != nil {
		return err
	}

	cacheKey := "timeline:user:" + strconv.FormatInt(userId, 10)

	var redisError error
	for _, tweetId := range tweetsId {
		date := time.Now().Unix()
		if err := cache.RedisClient.ZAdd(ctx, cacheKey, redis.Z{Score: float64(date), Member: tweetId}).Err(); err != nil {
			redisError = err
		}
	}

	if redisError != nil {
		return redisError
	}

	return nil
}
