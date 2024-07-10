package cache

import (
	"fmt"
	"strconv"
)

func GetUserFollowers(userId int64) ([]string, error) {
	cacheKey := "follower:user:" + strconv.FormatInt(userId, 10)

	cacheGet := RedisClient.SMembers(ctx, cacheKey)

	if cacheGet.Err() != nil || cacheGet.Val() == nil {
		return nil, cacheGet.Err()
	}

	return cacheGet.Val(), nil
}

func UpdateUserCacheTimeline(userId int64, tweetId int64, tweet string) error {
	fmt.Println(userId,tweetId,tweet)
	cacheKey := "tweets"

	cacheSet := RedisClient.HSet(ctx, cacheKey, tweetId, tweet)

	if cacheSet.Err() != nil {
		return cacheSet.Err()
	}

	followers, err := GetUserFollowers(userId)

	if err != nil {

		return err
	}

	fmt.Println("Fetched all followers: ",followers);

	go func() {
		for _, user := range followers {
			timelineCacheKey := "timeline:user:" + user

			fmt.Println("Current User is ",user)

			cacheUpdateTimeline := RedisClient.HSet(ctx, timelineCacheKey, tweetId, tweet)

			if cacheUpdateTimeline.Err() != nil {
				fmt.Printf("Error to update cache of userID %s", user)
			}
		}
	}()

	return nil
}
