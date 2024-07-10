package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"timeline/cmd/api/models"
	"timeline/cmd/cache"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func UpdateTweetCache(tweetId int64, tweet string) error {
	cacheKey := "tweets"
	err := cache.RedisClient.HSet(ctx, cacheKey, tweetId, tweet).Err()
	if err != nil {
		return err
	}

	log.Printf("[TWEET CACHE] Cache Updated with id %d ", tweetId)
	return nil
}

func UpdateUserCache(userId int64, tweetId int64) error {
	cacheKey := "user:tweets:" + strconv.FormatInt(userId, 10)

	date := time.Now().Unix()
	if err := cache.RedisClient.ZAdd(ctx, cacheKey, redis.Z{Score: float64(date), Member: tweetId}).Err(); err != nil {
		return err
	}

	log.Printf("[USER CACHE] User Cache Updated with id %d ", tweetId)
	return nil
}

func GetAllUserFollowers(userId int64) ([]string, error) {
	cacheKey := "follower:user:" + strconv.FormatInt(userId, 10)

	followers, err := cache.RedisClient.SMembers(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func UpdateUserTimeLine(userId string, tweetId int64) error {
	cacheKey := "timeline:user:" + userId

	date := time.Now().Unix()
	err := cache.RedisClient.ZAdd(ctx, cacheKey, redis.Z{Score: float64(date), Member: tweetId}).Err()

	if err != nil {
		fmt.Println(err)
		return err
	}

	log.Printf("[TIMELINE CACHE] Timeline Cache Updated for user with id %s with id %d ", userId, tweetId)
	return nil
}

func UpdateFollowersTimeLine(payload string) error {
	var Tweet models.Tweet
	err := json.Unmarshal([]byte(payload), &Tweet)
	if err != nil {
		return err
	}

	// Update Complete User Cache
	err = UpdateTweetCache(int64(Tweet.Id), payload)
	if err != nil {
		return err
	}

	// Update User Cache for uploaded tweet
	err = UpdateUserCache(int64(Tweet.UserId), int64(Tweet.Id))
	if err != nil {
		return err
	}

	followers, err := GetAllUserFollowers(int64(Tweet.UserId))

	log.Println(followers)
	if err != nil {
		return err
	}

	fmt.Println("All followers fetched successfully")

	// Update usertimeline
	go func() {
		for _, user := range followers {
			go UpdateUserTimeLine(user, int64(Tweet.Id))
		}
	}()

	return nil
}
