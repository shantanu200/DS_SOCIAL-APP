package controllers

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"user-relation/cmd/api/functions"
	"user-relation/cmd/api/handlers"
	"user-relation/cmd/api/models"
)

func FollowUserModel(c *fiber.Ctx) error {
	var relationBody *models.UserRelation

	if err := c.BodyParser(&relationBody); err != nil {
		return handlers.ErrorRouter(c, "Invalid Body Passed", err)
	}

	if relationBody.FolloweeId <= 0 || relationBody.FollowerId <= 0 {
		return handlers.ErrorRouter(c, "Invalid user ids passed", nil)
	}

	if relationBody.FolloweeId == relationBody.FollowerId {
		return handlers.ErrorRouter(c, "User cannot follow user personal account", nil)
	}

	err := functions.FollowUser(relationBody.FolloweeId, relationBody.FollowerId)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to follow user", err)
	}

	return handlers.SuccessRouter(c, "Followed User", relationBody)
}

func UnFollowUserModel(c *fiber.Ctx) error {
	var relationBody *models.UserRelation

	if err := c.BodyParser(&relationBody); err != nil {
		return handlers.ErrorRouter(c, "Invalid Body Passed", err)
	}

	if relationBody.FolloweeId <= 0 || relationBody.FollowerId <= 0 {
		return handlers.ErrorRouter(c, "Invalid user ids passed", nil)
	}

	if relationBody.FolloweeId == relationBody.FollowerId {
		return handlers.ErrorRouter(c, "User cannot un-follow user personal account", nil)
	}

	err := functions.UnFollowUser(relationBody.FolloweeId, relationBody.FollowerId)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to unfollow user", nil)
	}

	return handlers.SuccessRouter(c, "Unfollowed user", relationBody)
}

func UserFollowingModel(c *fiber.Ctx) error {
	id := c.Params("id")

	userId, _ := strconv.Atoi(id)

	var relations []models.UserRelation

	relations, err := functions.UserFollower(int64(userId))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch followers list", err)
	}

	return handlers.SuccessRouter(c, "Follower List found", relations)
}

func UserFollowerModel(c *fiber.Ctx) error {
	id := c.Params("id")

	userId, _ := strconv.Atoi(id)

	var relations []models.UserRelation

	relations, err := functions.UserFollowing(int64(userId))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch followers list", err)
	}

	return handlers.SuccessRouter(c, "Follower List found", relations)
}

func UserFollowingModelIds(c *fiber.Ctx) error {
	id := c.Params("id")

	userId, _ := strconv.Atoi(id)

	var relations []string

	relations, err := functions.UserFollowerIds(int64(userId))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch followers list of ids", err)
	}

	return handlers.SuccessRouter(c, "Follower List found of ids", relations)
}

func UserFollowerModelIds(c *fiber.Ctx) error {
	id := c.Params("id")

	userId, _ := strconv.Atoi(id)

	var relations []string

	relations, err := functions.UserFollowingIds(int64(userId))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch followers list of ids", err)
	}

	return handlers.SuccessRouter(c, "Follower List found of ids", relations)
}

func UserFollowingRecommendModel(c *fiber.Ctx) error {
	id := c.Params("id")

	userId, _ := strconv.Atoi(id)

	var relations *[]models.User

	relations, err := functions.UserFollowingRecommend(int64(userId))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch not  followed user list", err)
	}

	return handlers.SuccessRouter(c, "Not Followed User List found", relations)
}
