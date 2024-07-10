package routes

import (
	"user/cmd/api/controllers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func ServerRouter(app *fiber.App) {
	apiRouter := app.Group("/api")
	userRouter := apiRouter.Group("/user")

	userRouter.Post("/", controllers.CreateUserModel)
	userRouter.Post("/login", controllers.LoginUserModel)

	userRouter.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("modi_sarkar")},
	}))

	userRouter.Get("/details", controllers.GetUserByToken)
	userRouter.Patch("/details",controllers.UpdateUser)
}
