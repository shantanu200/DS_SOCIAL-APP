package functions

import (
	"fmt"
	"time"
	"tweet/cmd/api/models"
	"tweet/cmd/cache"
	"tweet/cmd/database"
	rabbitmq "tweet/cmd/rabbitMQ"

	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

func CreateTweet(tweet models.Tweet) (*models.Tweet, error) {
	db := database.DB

	tweetBody := models.Tweet{
		UserId:       tweet.UserId,
		CreationTime: time.Now(),
		Content:      tweet.Content,
		TotalLikes:   0,
		MediaFiles:   tweet.MediaFiles,
	}

	err := db.Create(&tweetBody).Error
	if err != nil {
		return nil, err
	}

	var tweetObj *models.Tweet

	err = db.Model(models.Tweet{}).Preload("User").First(&tweetObj, tweetBody.Id).Error

	if err != nil {
		return nil, err
	}

	go func() {
		jsonMarshal, err := json.Marshal(tweetObj)

		if err != nil {
			fmt.Println("JSON Marshal error: ", err.Error())
			return
		}

		err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME, "FANOUT", string(jsonMarshal))

		if err != nil {
			fmt.Println("RabbitMQ error: ", err.Error())
			return
		}
	}()

	return &tweetBody, nil
}

func GetAllTweets() (*[]models.Tweet, error) {
	db := database.DB

	var tweets []models.Tweet

	cacheKey := "tweets"

	cacheGet := cache.RedisClient.HGetAll(ctx, cacheKey)

	if cacheGet.Err() == nil && len(cacheGet.Val()) > 0 {
		start := time.Now()
		for _, tweet := range cacheGet.Val() {
			var tweetObj models.Tweet
			err := json.Unmarshal([]byte(tweet), &tweetObj)

			if err != nil {
				return nil, err
			}

			tweets = append(tweets, tweetObj)
		}

		fmt.Println(time.Since(start))
	} else {
		fmt.Println("Cache MISS")
		err := db.Model(models.Tweet{}).Preload("User").Order("created_at desc").Find(&tweets).Error
		if err != nil {
			return nil, err
		}

		if len(tweets) > 0 {
			for _, tweet := range tweets {
				jsonMarshal, err := json.Marshal(tweet)

				if err != nil {
					return nil, err
				}

				tweetPush := cache.RedisClient.HSet(ctx, cacheKey, tweet.Id, string(jsonMarshal))

				if tweetPush.Err() != nil {
					return nil, tweetPush.Err()
				}
			}
		}
	}

	return &tweets, nil
}

func GetLikeUserList(tweetId int) (*[]models.FavouriteTweet, error) {
	db := database.DB

	var users *[]models.FavouriteTweet

	err := db.Model(models.FavouriteTweet{}).Select("user_id").Where("tweet_id = ?", tweetId).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func FavTweet(userId int64, tweetId int64, isLike bool) error {
	db := database.DB

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if isLike {
		_, err := LikeTweet(userId, tweetId)
		if err != nil {
			tx.Rollback()
			return err
		}

		likeErr := db.Model(&models.Tweet{}).Where("id = ?", tweetId).Update("total_likes", gorm.Expr("total_likes + ?", 1)).Error

		if likeErr != nil {
			tx.Rollback()
			return likeErr
		}
	} else {
		_, err := UnlikeTweet(userId, tweetId)
		if err != nil {
			tx.Rollback()
			return err
		}

		unlikeErr := db.Model(&models.Tweet{}).Where("id = ?", tweetId).Update("total_likes", gorm.Expr("total_likes - ?", 1)).Error

		if unlikeErr != nil {
			tx.Rollback()
			return unlikeErr
		}
	}

	return tx.Commit().Error
}

func GetSingleTweet(tweetId string) (*models.Tweet, error) {
	cacheKey := "tweets"

	result, err := cache.RedisClient.HGet(ctx, cacheKey, tweetId).Result()

	if err != nil {
		return nil, err
	}

	var tweet *models.Tweet

	err = json.Unmarshal([]byte(result), &tweet)

	if err != nil {
		return nil, err
	}

	return tweet, nil
}
