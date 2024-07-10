package routes

import (
	"tweet/cmd/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func ServerRouter(app *fiber.App) {
	apiRouter := app.Group("/api")
	tweetRouter := apiRouter.Group("/tweet")

	tweetRouter.Post("/", controllers.CreateTweetModel)
	tweetRouter.Get("/", controllers.GetAllTweets)
	tweetRouter.Get("/details/:id", controllers.GetSingleTweet)
	tweetRouter.Patch("/action", controllers.UserActionModel)
	tweetRouter.Get("/action/:id", controllers.GetAllLikedUserForTweet)
	tweetRouter.Get("/user/action/:id", controllers.UserLikedModel)
	tweetRouter.Post("/reply", controllers.CreateReplyModel)
	tweetRouter.Get("/reply/:id", controllers.GetReplyOnTweet)
	tweetRouter.Patch("/reply/action", controllers.UserReplyActionModel)
	tweetRouter.Get("/reply/user/action/:id", controllers.UserLikeReplyIdModel)
	tweetRouter.Post("/thread",controllers.CreateThreadModel);
	tweetRouter.Get("/thread/:id",controllers.GetAllThreadRepliesModel)
}
