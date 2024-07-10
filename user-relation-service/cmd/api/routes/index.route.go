package routes

import (
	"github.com/gofiber/fiber/v2"
	"user-relation/cmd/api/controllers"
)

func ServerRouter(app *fiber.App) {
	apiRouter := app.Group("/api")
	relationRouter := apiRouter.Group("/relation")
	relationRouter.Patch("/follow", controllers.FollowUserModel)
	relationRouter.Patch("/unFollow", controllers.UnFollowUserModel)
	relationRouter.Get("/following/:id", controllers.UserFollowingModel)
	relationRouter.Get("/follower/:id", controllers.UserFollowerModel)
	relationRouter.Get("/follower_id/:id", controllers.UserFollowerModelIds)
	relationRouter.Get("/following_id/:id", controllers.UserFollowingModelIds)
	relationRouter.Get("/recommend/:id", controllers.UserFollowingRecommendModel)
}
