package main

import (
	"user/cmd/api/routes"
	"user/cmd/database"
	rabbitmq "user/cmd/rabbitMQ"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.ConnectPostgres()
	rabbitmq.Connect()
	rabbitmq.SetupConsumeQueue()
	userApp := fiber.New()
	userApp.Use(healthcheck.New())
	userApp.Use(compress.New())
	userApp.Use(cors.New())
	userApp.Use(logger.New())
	userApp.Get("/service", func(c *fiber.Ctx) error {
		return c.SendString("Hello from User Service")
	})
	routes.ServerRouter(userApp)
	userApp.Listen(":80")
}
