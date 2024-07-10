package cachefunction

import (
	"strconv"
	"time"
	"timeline/cmd/cache"

	"github.com/redis/go-redis/v9"
)

func UpdateUserLikeReplyCache(replyId int64, userId int64) error {
	userReplyLikeKey := "userReplyLike:" + strconv.FormatInt(int64(userId), 10)

	date := time.Now().Unix()

	if err := cache.RedisClient.ZAdd(ctx, userReplyLikeKey, redis.Z{Score: float64(date), Member: replyId}).Err(); err != nil {
		return err
	}

	return nil
}

func UpdateUserDisLikeReplyCache(replyId int64, userId int64) error {
	userReplyLikeKey := "userReplyLike:" + strconv.FormatInt(int64(userId), 10)

	if err := cache.RedisClient.ZRem(ctx, userReplyLikeKey, replyId).Err(); err != nil {
		return err
	}

	return nil
}
