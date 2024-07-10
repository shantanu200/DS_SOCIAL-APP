package event

import (
	"log"
	"strconv"
	"timeline/cmd/cache"
	cachefunction "timeline/cmd/cache/cache_function"

	"github.com/goccy/go-json"
)

type FollowPayload struct {
	IsFollow   bool  `json:"isFollow"`
	FolloweeId int64 `json:"followeeId"`
	FollowerId int64 `json:"followerId"`
}

func UpdateFollowingFeed(followingUserId int64, userId int64) error {
	cacheKey := "user:tweets:" + strconv.FormatInt(userId, 10)

	tweets, err := cache.RedisClient.ZRevRange(ctx, cacheKey, 0, 3).Result()

	if err != nil {
		return err
	}

	_followingUserId := strconv.FormatInt(followingUserId, 10)

	for _, tweetId := range tweets {
		_tweetId, _ := strconv.Atoi(tweetId)
		if err := UpdateUserTimeLine(_followingUserId, int64(_tweetId)); err != nil {
			return err
		}
	}

	return nil
}

func AddFollower(cacheKey string, followeeId int64, userId int64) error {
	err := cache.RedisClient.SAdd(ctx, cacheKey, userId).Err()

	if err != nil {
		return err
	}

	if err := UpdateFollowingFeed(followeeId,userId); err != nil {
		return err
	}

	log.Println("Follower is added to follower list")
	return nil
}

func RemoveFollower(cacheKey string, removeUserId int64, userId int64) error {
	err := cache.RedisClient.SRem(ctx, cacheKey, removeUserId).Err()

	if err != nil {
		return err
	}

	err = cachefunction.RemoveTweetsFromTimeline(removeUserId, userId)

	if err != nil {
		return err
	}

	log.Println("Follower is remove from follower list")
	return nil
}

func UpdateFollowerUserList(payload string) error {
	var followPayload FollowPayload

	if err := json.Unmarshal([]byte(payload), &followPayload); err != nil {
		return err
	}

	cacheKey := "following:user:" + strconv.FormatInt(followPayload.FolloweeId, 10)
	if followPayload.IsFollow {
		if err := AddFollower(cacheKey, followPayload.FolloweeId, followPayload.FollowerId); err != nil {
			return err
		}
	} else {
		if err := RemoveFollower(cacheKey, followPayload.FollowerId, followPayload.FolloweeId); err != nil {
			return err
		}
	}

	return nil
}
