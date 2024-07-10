package main

import (
	"user-relation/cmd/api/routes"
	"user-relation/cmd/cache"
	"user-relation/cmd/database"
	rabbitmq "user-relation/cmd/rabbitMQ"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.ConnectPostgres()
	cache.ConnectCache()
	rabbitmq.Connect()
	rabbitmq.SetupConsumeQueue()

	userRelation := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	userRelation.Use(healthcheck.New())

	userRelation.Use(logger.New())

	userRelation.Use(cors.New())

	routes.ServerRouter(userRelation)

	userRelation.Listen(":80")
}
