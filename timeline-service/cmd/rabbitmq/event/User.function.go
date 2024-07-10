package event

import (
	cachefunction "timeline/cmd/cache/cache_function"
	dbhandler "timeline/cmd/database/db.handler"

	"github.com/goccy/go-json"
)

type UpdateUserPayload struct {
	UserId int64 `json:"userId"`
}

func UpdateUserTweets(payload string) error {
	var updateUserPaylod UpdateUserPayload

	err := json.Unmarshal([]byte(payload), &updateUserPaylod)
	if err != nil {
		return err
	}

	tweetIds, err := cachefunction.GetAllUserTweets(updateUserPaylod.UserId)

	if err != nil {
		return err
	}

	tweets, err := dbhandler.GetAllTweetsForUser(tweetIds)

	if err != nil {
		return err
	}
    
	err = cachefunction.UpdateTweets(tweets)

	if err != nil {
		return err
	}

	return nil
}
