package main

import (
	"tweet/cmd/api/routes"
	"tweet/cmd/cache"
	"tweet/cmd/database"
	rabbitmq "tweet/cmd/rabbitMQ"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.ConnectPostgres()
	cache.ConnectCache()
	rabbitmq.Connect()
	rabbitmq.SetupConsumeQueue()

	tweetApp := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   10 * 1024 * 1024,
	})

	tweetApp.Use(healthcheck.New())
	tweetApp.Use(cors.New())
	tweetApp.Use(compress.New())
	tweetApp.Use(logger.New())

	tweetApp.Get("/service", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Tweet Service")
	})

	routes.ServerRouter(tweetApp)

	tweetApp.Listen(":80")
}
