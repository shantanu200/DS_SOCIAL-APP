package controllers

import (
	"strconv"
	"timeline/cmd/api/functions"
	"timeline/cmd/api/handlers"
	cachefunction "timeline/cmd/cache/cache_function"

	"github.com/gofiber/fiber/v2"
)

func GetUserTimeLineModel(c *fiber.Ctx) error {
	id := c.Params("id")
	Qpage := c.Query("page", "1")
	QSize := c.Query("size", "10")

	UID, _ := strconv.Atoi(id)
	UPage, _ := strconv.Atoi(Qpage)
	USize, _ := strconv.Atoi(QSize)

	results, err := functions.GetUserTimeLine(int64(UID), int64(UPage), int64(USize))
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(fiber.Map{"data": results})
}

func GetUserHomeTimeline(c *fiber.Ctx) error {
	id := c.Params("id")
	Qpage := c.Query("page", "1")
	QSize := c.Query("size", "10")

	UID, _ := strconv.Atoi(id)
	UPage, _ := strconv.Atoi(Qpage)
	USize, _ := strconv.Atoi(QSize)

	results, err := cachefunction.GetHomeTimelineCache(int64(UID), int64(UPage), int64(USize))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch results")
	}

	return handlers.SuccessRouter(c, "User home timeline fetched", results)
}

func GetUserLikeTweetsModel(c *fiber.Ctx) error {
	id := c.Params("id")

	UID, _ := strconv.Atoi(id)

	results, err := cachefunction.GetUserLikedTweets(int64(UID))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch results")
	}

	return handlers.SuccessRouter(c, "User home timeline fetched", results)
}

func GetUserPostTweetsModel(c *fiber.Ctx) error {
	id := c.Params("id")

	UID, _ := strconv.Atoi(id)

	results, err := cachefunction.GetUserPostTweets(int64(UID))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch results")
	}

	return handlers.SuccessRouter(c, "User home timeline fetched", results)
}
