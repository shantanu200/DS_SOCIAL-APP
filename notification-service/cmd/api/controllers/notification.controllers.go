package controllers

import (
	"notification/cmd/api/handlers"
	"notification/cmd/socket/functions"
	"notification/cmd/socket/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AssignNotificationToUserModel(c *fiber.Ctx) error {
	var notifications []models.Notification

	if err := c.BodyParser(&notifications); err != nil {
		return handlers.ErrorRouter(c, "Invalid Body passed", err)
	}

	if len(notifications) == 0 {
		return handlers.ErrorRouter(c, "Empty notification arrray", nil)
	}

	err := functions.CreateNotifications(notifications)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to insert multiple notification", err)
	}

	return handlers.SuccessRouter(c, "Notification added", notifications)
}

type ReadPayload struct {
	UserId          int64   `json:"userId"`
	NotificationIds []int64 `json:"notificationIds"`
}

func ReadUserNotifcationModel(c *fiber.Ctx) error {
	var readNotificationBody ReadPayload

	if err := c.BodyParser(&readNotificationBody); err != nil {

		return handlers.ErrorRouter(c, "Invalid Body passed", err)
	}

	if readNotificationBody.UserId <= 0 {
		return handlers.ErrorRouter(c, "UserId required for process", nil)
	}

	if len(readNotificationBody.NotificationIds) == 0 {
		return handlers.ErrorRouter(c, "Empty notification arrray", nil)
	}

	err := functions.ReadNotifications(readNotificationBody.UserId, readNotificationBody.NotificationIds)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to update notification", err)
	}

	return handlers.SuccessRouter(c, "Notification updated", readNotificationBody)
}

func NotificationForUserModel(c *fiber.Ctx) error {
	id := c.Params("id")

	userId, _ := strconv.Atoi(id)

	results, err := functions.GetAllUserNotifications(int64(userId))

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to fetch notification", err)
	}

	return handlers.SuccessRouter(c, "All notification for user", results)
}
