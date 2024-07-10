package event

import (
	"fmt"
	"strconv"
	"timeline/cmd/api/models"
	"timeline/cmd/cache"

	"github.com/goccy/go-json"
)

type LikePayload struct {
	UserId  int64 `json:"userId"`
	TweetId int64 `json:"tweetId"`
	IsLike  bool  `json:"isLike"`
}

func UpdateTweet(tweetCacheKey string, tweetId int64, tweetString string) error {
	err := cache.RedisClient.HSet(ctx, tweetCacheKey, tweetId, tweetString).Err()

	if err != nil {
		return err
	}

	fmt.Println("Tweet likes updated")

	return nil
}

func UpdateLikeTweet(userCacheKey string, tweetId int64) error {
	err := cache.RedisClient.SAdd(ctx, userCacheKey, tweetId).Err()

	if err != nil {
		return err
	}

	fmt.Println("Tweet added to user like")

	return nil
}

func UpdateDislikeTweet(userCacheKey string, tweetId int64) error {
	err := cache.RedisClient.SRem(ctx, userCacheKey, tweetId).Err()

	if err != nil {
		return err
	}

	fmt.Println("Tweet deleted from user like")
	return nil
}

func UpdateLikeStatus(payload string) error {
	var likePayload *LikePayload
	err := json.Unmarshal([]byte(payload), &likePayload)

	if err != nil {
		return err
	}

	tweetCacheKey := "tweets"
	userLikeCacheKey := "userLike:" + strconv.FormatInt(int64(likePayload.UserId), 10)
	sTweetId := strconv.FormatInt(int64(likePayload.TweetId), 10)

	tweet, err := cache.RedisClient.HGet(ctx, tweetCacheKey, sTweetId).Result()

	if err != nil {
		return err
	}

	var CacheTweet *models.Tweet

	err = json.Unmarshal([]byte(tweet), &CacheTweet)

	if err != nil {
		return err
	}

	if likePayload.IsLike {
		if CacheTweet.TotalLikes >= 0 {
			CacheTweet.TotalLikes++
		}

		jsonMarshal, err := json.Marshal(CacheTweet)

		if err != nil {
			return err
		}

		err = UpdateTweet(tweetCacheKey, int64(CacheTweet.Id), string(jsonMarshal))

		if err != nil {
			return err
		}

		err = UpdateLikeTweet(userLikeCacheKey, likePayload.TweetId)

		if err != nil {
			return err
		}
	} else {
		if CacheTweet.TotalLikes > 0 {
			CacheTweet.TotalLikes--
		}

		jsonMarshal, err := json.Marshal(CacheTweet)

		if err != nil {
			return err
		}

		err = UpdateTweet(tweetCacheKey, int64(CacheTweet.Id), string(jsonMarshal))

		if err != nil {
			return err
		}

		err = UpdateDislikeTweet(userLikeCacheKey, likePayload.TweetId)

		if err != nil {
			return err
		}
	}

	return nil
}
