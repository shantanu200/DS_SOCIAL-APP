package functions

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"user-relation/cmd/api/models"
	"user-relation/cmd/cache"
	"user-relation/cmd/database"
	rabbitmq "user-relation/cmd/rabbitMQ"

	"github.com/goccy/go-json"
)

var ctx context.Context = context.Background()

func FollowUser(followeeId int64, followerId int64) error {
	db := database.DB

	userRelation := models.UserRelation{
		FolloweeId:   followeeId,
		FollowerId:   followerId,
		CreationTime: time.Now(),
	}

	err := db.Model(&models.UserRelation{}).Create(&userRelation).Error

	if err != nil {
		if exists := strings.Contains(err.Error(), "duplicate key"); exists {
			return nil
		} else {
			return err
		}
	}
    
	var followPayload = rabbitmq.FollowerPayload{
		FolloweeId: followeeId,
		FollowerId: followerId,
		IsFollow: true,
	}

	jsonMarshal,err := json.Marshal(followPayload)

	if err != nil {
		return err
	}

	err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME,"FOLLOWER",string(jsonMarshal))

	if err != nil {
		return err
	}
  
	return nil
}

func UnFollowUser(followeeId int64, followerId int64) error {
	db := database.DB

	err := db.Where("followee_id = ? AND follower_id = ?", followeeId, followerId).Delete(&models.UserRelation{}).Error

	if err != nil {
		fmt.Println(err)
		return err
	} 
	var followPayload = rabbitmq.FollowerPayload{
		FolloweeId: followeeId,
		FollowerId: followerId,
		IsFollow: false,
	}

	jsonMarshal,err := json.Marshal(followPayload)

	if err != nil {
		return err
	}

	err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME,"FOLLOWER",string(jsonMarshal))

	if err != nil {
		return err
	}

	return nil
}

// Returns users followers list (which are user following)
func UserFollower(userId int64) ([]models.UserRelation, error) {
	db := database.DB
	cacheKey := "following:user:" + strconv.FormatInt(userId, 10)

	var followUserModels []models.UserRelation

	err := db.Preload("Follower").Model(&models.UserRelation{}).Where("followee_id = ?", userId).Find(&followUserModels).Error
	if err != nil {
		return nil, err
	}

	cacheGet := cache.RedisClient.SMembers(ctx, cacheKey)

	if cacheGet.Err() != nil || cacheGet.Val() == nil {
for _, users := range followUserModels {
			cachePush := cache.RedisClient.SAdd(ctx, cacheKey, users.FollowerId)

			if cachePush.Err() != nil {
				return nil, cachePush.Err()
			}
		}
	}		

	return followUserModels, nil
}
func UserFollowerIds(userId int64) ([]string, error) {
	db := database.DB
	cacheKey := "following:user:" + strconv.FormatInt(userId, 10)
	var followUserModels []models.UserRelation

	cacheGet := cache.RedisClient.SMembers(ctx, cacheKey)

	if cacheGet.Err() == nil && len(cacheGet.Val()) > 0 {
		return cacheGet.Val(), nil
	} else {
		var userIds []string
		err := db.Preload("Follower").Model(&models.UserRelation{}).Where("followee_id = ?", userId).Find(&followUserModels).Error
		if err != nil {
			return nil, err
		}
		for _, users := range followUserModels {
			cachePush := cache.RedisClient.SAdd(ctx, cacheKey, users.FollowerId)

			userIds = append(userIds, strconv.FormatInt(users.FollowerId, 10))

			if cachePush.Err() != nil {
				return nil, cachePush.Err()
			}
		}

		return userIds, nil
	}
}

func UserFollowingIds(userId int64) ([]string, error) {
	db := database.DB
	cacheKey := "follower:user:" + strconv.FormatInt(userId, 10)
	var followUserModels []models.UserRelation

	cacheGet := cache.RedisClient.SMembers(ctx, cacheKey)

	if cacheGet.Err() == nil && len(cacheGet.Val()) > 0 {
		return cacheGet.Val(), nil
	} else {
		var userIds []string
		err := db.Preload("Followee").Model(&models.UserRelation{}).Where("follower_id = ?", userId).Find(&followUserModels).Error
		if err != nil {
			return nil, err
		}
		for _, users := range followUserModels {
			cachePush := cache.RedisClient.SAdd(ctx, cacheKey, users.FolloweeId)

			userIds = append(userIds, strconv.FormatInt(users.FolloweeId, 10))

			if cachePush.Err() != nil {
				return nil, cachePush.Err()
			}
		}

		return userIds, nil
	}
}

// Returns users followee list (which are following user)
func UserFollowing(userId int64) ([]models.UserRelation, error) {
	db := database.DB
	cacheKey := "follower:user:" + strconv.FormatInt(userId, 10)
	var followUserModels []models.UserRelation

	cacheGet := cache.RedisClient.SMembers(ctx, cacheKey)
	err := db.Preload("Followee").Model(&models.UserRelation{}).Where("follower_id = ?", userId).Find(&followUserModels).Error
	if err != nil {
		return nil, err
	}
	if cacheGet.Err() != nil || cacheGet.Val() == nil {
		for _, users := range followUserModels {
			cachePush := cache.RedisClient.SAdd(ctx, cacheKey, users.FolloweeId)

			if cachePush.Err() != nil {
				return nil, cachePush.Err()
			}
		}
	}

	return followUserModels, nil
}

func UserFollowingRecommend(userId int64) (*[]models.User, error) {
	db := database.DB
	cacheKey := "following:user:" + strconv.FormatInt(userId, 10)
	var users *[]models.User

	cacheGet := cache.RedisClient.SMembers(ctx, cacheKey)

	if cacheGet.Err() == nil && len(cacheGet.Val()) > 0 {
		err := db.Model(&models.User{}).Where("id NOT IN (?) AND id != ?", cacheGet.Val(), userId).Find(&users).Limit(3).Error

		if err != nil {
			return nil, err
		}
	} else {
		subQuery := db.Model(&models.UserRelation{}).Select("follower_id").Where("followee_id = ?", userId)

		err := db.Model(&models.User{}).Where("id NOT IN (?) AND id != ?", subQuery, userId).Find(&users).Limit(3).Error

		if err != nil {
			return nil, err
		}
	}

	return users, nil
}
