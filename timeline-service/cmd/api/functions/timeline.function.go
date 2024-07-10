package functions

import (
	"timeline/cmd/api/models"
	cachefunction "timeline/cmd/cache/cache_function"
)

func GetUserTimeLine(userId int64,page int64,size int64) ([]models.Tweet, error) {
	tweetId, err := cachefunction.GetUserTimelineCache(userId,page,size)
	if err != nil {
		return nil, err
	}

	tweets, err := cachefunction.GetTweets(tweetId)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}
