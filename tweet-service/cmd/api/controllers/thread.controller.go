package controllers

import (
	"strconv"
	"tweet/cmd/api/functions"
	"tweet/cmd/api/handlers"
	"tweet/cmd/api/models"
	"tweet/cmd/s3"

	"github.com/gofiber/fiber/v2"
)

func CreateThreadModel(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return handlers.ErrorRouter(c, "Invalid Form Data recieved", nil)
	}

	replyId, err := strconv.Atoi(form.Value["replyId"][0])

	if err != nil || replyId <= 0 {
		return handlers.ErrorRouter(c, "Please insert valid replyId", nil)
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

	var threadBody = models.Thread{
		UserId:     userId,
		ReplyId:    replyId,
		Content:    content,
		MediaFiles: mediaUrls,
		TweetId:    tweetId,
	}

	tweet, err := functions.CreateThread(&threadBody)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable post a reply on tweet", err)
	}

	return handlers.SuccessRouter(c, "Reply posted on tweet", tweet)
}

func GetAllThreadRepliesModel(c *fiber.Ctx) error {
	id := c.Params("id")
	QPage := c.QueryInt("page")
	QSize := c.QueryInt("size")
	UID, _ := strconv.Atoi(id)

	results, err := functions.GetAllThreadReplies(int64(UID), int64(QPage), int64(QSize))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch a thread of tweet", err)
	}

	return handlers.SuccessRouter(c, "Thread fetched for a tweet", results)
}
