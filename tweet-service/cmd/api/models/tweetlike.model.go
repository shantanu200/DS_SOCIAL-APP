package models

import "time"

type FavouriteTweet struct {
	Id           int64     `json:"id" gorm:"primaryKey"`
	UserId       int64     `json:"userid" gorm:"uniqueIndex:favourite_idx,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreationTime time.Time `json:"createdAt"`
	TweetId      int64     `json:"tweetid" gorm:"uniqueindex:favourite_idx,constraint:onupdate:cascade,ondelete:set null;"`
}
