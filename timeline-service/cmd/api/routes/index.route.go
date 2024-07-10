package routes

import (
	"timeline/cmd/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func ServerRouter(app *fiber.App) {
	api := app.Group("/api")
	timeLineRouter := api.Group("/timeline")

	timeLineRouter.Get("/following/:id", controllers.GetUserTimeLineModel)
	timeLineRouter.Get("/home/:id", controllers.GetUserHomeTimeline)
	timeLineRouter.Get("/like/:id", controllers.GetUserLikeTweetsModel)
	timeLineRouter.Get("/posts/:id", controllers.GetUserPostTweetsModel)
}
