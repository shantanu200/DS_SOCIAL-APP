package event

import (
	cachefunction "timeline/cmd/cache/cache_function"

	"github.com/goccy/go-json"
)

type ReplyLikePayload struct {
	UserId  int64 `json:"userId"`
	ReplyId int64 `json:"replyId"`
	IsLike  bool  `json:"isLike"`
}

func UpdateUserReplyLike(payload string) error {
	var replyPayload *ReplyLikePayload

	if err := json.Unmarshal([]byte(payload), &replyPayload); err != nil {
		return err
	}

	if replyPayload.IsLike {
		if err := cachefunction.UpdateUserLikeReplyCache(replyPayload.ReplyId, replyPayload.UserId); err != nil {
			return err
		}
	} else {
		if err := cachefunction.UpdateUserDisLikeReplyCache(replyPayload.ReplyId, replyPayload.UserId); err != nil {
			return err
		}
	}

	return nil
}
