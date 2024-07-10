package controllers

import (
	"strconv"
	"tweet/cmd/api/functions"
	"tweet/cmd/api/handlers"
	"tweet/cmd/api/models"
	"tweet/cmd/s3"

	"github.com/gofiber/fiber/v2"
)

func CreateReplyModel(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return handlers.ErrorRouter(c, "Invalid Form Data recieved", nil)
	}

	userId, err := strconv.Atoi(form.Value["userId"][0])

	if err != nil || userId <= 0 {
		return handlers.ErrorRouter(c, "Please insert valid userId", nil)
	}

	tweetId, err := strconv.Atoi(form.Value["tweetId"][0])

	if err != nil || tweetId <= 0 {
		return handlers.ErrorRouter(c, "Please insert valid tweetId", nil)
	}

	content := form.Value["content"][0]

	if len(content) > 140 {
		return handlers.ErrorRouter(c, "Tweet content must be less than 140 chars", nil)
	}

	if content == "" {
		return handlers.ErrorRouter(c, "Empty tweet is not allowed", nil)
	}

	mediaFiles := form.File["mediaFiles"]

	if len(mediaFiles) > 5 {
		return handlers.ErrorRouter(c, "Max 5 files can be uploaded for the post", nil)
	}

	var mediaUrls []string

	if len(mediaFiles) > 0 {
		results, err := s3.FileUploader(mediaFiles)

		if err != nil {
			return handlers.ErrorRouter(c, "Unable to upload media files on s3 bucket", nil)
		}

		mediaUrls = results
	}

	var replyBody = models.Reply{
		Content:    content,
		UserId:     userId,
		MediaFiles: mediaUrls,
		TweetId:    tweetId,
	}

	tweet, err := functions.CreateReply(replyBody)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable post a reply on tweet", err)
	}

	return handlers.SuccessRouter(c, "Reply posted on tweet", tweet)
}

func GetReplyOnTweet(c *fiber.Ctx) error {
	id := c.Params("id")
	qPage := c.Query("page", "1")
	qSize := c.Query("size", "10")

	UID, _ := strconv.Atoi(id)
	UPage, _ := strconv.Atoi(qPage)
	USize, _ := strconv.Atoi(qSize)

	replies, err := functions.GetAllReplies(int64(UID), int64(UPage), int64(USize))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable get replies on tweet", err)
	}

	return handlers.SuccessRouter(c, "Reply posted on tweet", replies)
}

func UserReplyActionModel(c *fiber.Ctx) error {
	var replyBody map[string]interface{}

	if err := c.BodyParser(&replyBody); err != nil {
		return handlers.ErrorRouter(c, "Invalid Body Passed", err)
	}

	if _, exists := replyBody["userId"]; !exists {
		return handlers.ErrorRouter(c, "Please insert valid userId", nil)
	}

	if _, exists := replyBody["replyId"]; !exists {
		return handlers.ErrorRouter(c, "Please insert valid replyId", nil)
	}

	err := functions.LikeReplyAction(int64(replyBody["replyId"].(float64)), int64(replyBody["userId"].(float64)), replyBody["isLike"].(bool))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to like tweet", err)
	}

	return handlers.SuccessRouter(c, "Tweet action executed", replyBody)
}

func UserLikeReplyIdModel(c *fiber.Ctx) error {
	id := c.Params("id")

	UID, _ := strconv.Atoi(id)

	results, err := functions.UserLikeRepliesIds(int64(UID))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to retrun like tweet ids", err)
	}

	return handlers.SuccessRouter(c, "Tweet action executed", results)
}
