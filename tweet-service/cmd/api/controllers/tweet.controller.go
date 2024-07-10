package controllers

import (
	"strconv"
	"tweet/cmd/api/functions"
	"tweet/cmd/api/handlers"
	"tweet/cmd/api/models"
	"tweet/cmd/s3"

	"github.com/gofiber/fiber/v2"
)

func CreateTweetModel(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return handlers.ErrorRouter(c, "Invalid Form Data recieved", nil)
	}

	userId, err := strconv.Atoi(form.Value["userId"][0])

	if err != nil || userId <= 0 {
		return handlers.ErrorRouter(c, "Please insert valid userId", nil)
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

	var tweetBody = models.Tweet{
		Content:    content,
		UserId:     userId,
		MediaFiles: mediaUrls,
	}

	tweet, err := functions.CreateTweet(tweetBody)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable post a tweet", err)
	}

	return handlers.SuccessRouter(c, "Tweet posted", tweet)
}

func GetAllTweets(c *fiber.Ctx) error {
	tweets, err := functions.GetAllTweets()

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch tweets", err)
	}

	return handlers.SuccessRouter(c, "All tweets", &tweets)
}

func GetSingleTweet(c *fiber.Ctx) error {
	id := c.Params("id")
	tweet, err := functions.GetSingleTweet(id)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch tweets", err)
	}

	return handlers.SuccessRouter(c, "All tweets", &tweet)
}

func GetAllLikedUserForTweet(c *fiber.Ctx) error {
	id := c.Params("id")

	I_id, _ := strconv.Atoi(id)

	users, err := functions.GetLikeUserList(I_id)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to liked user tweets", err)
	}

	return handlers.SuccessRouter(c, "All tweets", &users)
}

func UserLikedModel(c *fiber.Ctx) error {
	id := c.Params("id")

	I_id, _ := strconv.Atoi(id)

	users, err := functions.UserLikeTweets(int64(I_id))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to get all user liked post", err)
	}

	return handlers.SuccessRouter(c, "All user liked post", &users)
}

func UserActionModel(c *fiber.Ctx) error {
	var tweetBody map[string]interface{}

	if err := c.BodyParser(&tweetBody); err != nil {
		return handlers.ErrorRouter(c, "Invalid Body Passed", err)
	}

	if _, exists := tweetBody["userId"]; !exists {
		return handlers.ErrorRouter(c, "Please insert valid userId", nil)
	}

	if _, exists := tweetBody["tweetId"]; !exists {
		return handlers.ErrorRouter(c, "Please insert valid tweetId", nil)
	}

	err := functions.FavTweet(int64(tweetBody["userId"].(float64)), int64(tweetBody["tweetId"].(float64)), tweetBody["isLike"].(bool))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to like tweet", err)
	}

	return handlers.SuccessRouter(c, "Tweet action executed", tweetBody)
}
