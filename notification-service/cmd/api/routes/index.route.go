package routes

import (
	"notification/cmd/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func ServerRouter(app *fiber.App) {
	apiRouter := app.Group("/api")
	notificationRouter := apiRouter.Group("/notification")

	notificationRouter.Post("/", controllers.AssignNotificationToUserModel)
	notificationRouter.Patch("/read", controllers.ReadUserNotifcationModel)
	notificationRouter.Get("/all/:id",controllers.NotificationForUserModel)
}
