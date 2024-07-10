package models

import (
	"time"

	"github.com/lib/pq"
)

type Reply struct {
	Id           int            `json:"id" gorm:"primaryKey"`
	TweetId      int            `json:"tweetId" form:"tweetId"`
	Tweet        Tweet          `json:"tweet" gorm:"foreignKey:TweetId"`
	UserId       int            `json:"userId" form:"userId"`
	User         User           `json:"user" gorm:"foreignKey:UserId"`
	CreationTime time.Time      `json:"createdAt"`
	Content      string         `json:"content" form:"content"`
	TotalLikes   int            `json:"totalLikes"`
	MediaFiles   pq.StringArray `json:"mediaFiles" gorm:"type:text[]" form:"mediaFiles"`
	TotalReplies int            `json:"totalReplies"`
}

type Thread struct {
	Id           int            `json:"id" gorm:"primaryKey"`
	ThreadId     int            `json:"threadId"`
	ReplyId      int            `json:"ReplyId" form:"ReplyId"`
	Reply        Reply          `json:"reply" gorm:"foreignKey:ReplyId"`
	TweetId      int            `json:"tweetId" form:"tweetId"`
	Tweet        Tweet          `json:"tweet" gorm:"foreignKey:TweetId"`
	UserId       int            `json:"userId" form:"userId"`
	User         User           `json:"user" gorm:"foreignKey:UserId"`
	CreationTime time.Time      `json:"createdAt"`
	Content      string         `json:"content" form:"content"`
	TotalLikes   int            `json:"totalLikes"`
	MediaFiles   pq.StringArray `json:"mediaFiles" gorm:"type:text[]" form:"mediaFiles"`
}

type FavouriteReply struct {
	Id           int       `json:"id" gorm:"primaryKey"`
	ReplyId      int       `json:"replyid" gorm:"uniqueindex:favourite_reply_idx,constraint:onupdate:cascade,ondelete:set null;"`
	UserId       int       `json:"userid" gorm:"uniqueIndex:favourite_reply_idx,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreationTime time.Time `json:"createdAt"`
}
