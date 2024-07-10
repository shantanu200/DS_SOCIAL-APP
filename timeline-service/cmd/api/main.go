package main

import (
	"log"
	"timeline/cmd/api/routes"
	"timeline/cmd/cache"
	"timeline/cmd/database"
	"timeline/cmd/rabbitmq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cache.ConnectCache()
	database.ConnectPostgres()
	rabbitConn, err := rabbitmq.Connect()
	if err != nil {
		log.Println(err)
	}
	defer rabbitConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := rabbitmq.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	timeLine := fiber.New()

	timeLine.Use(healthcheck.New())
	timeLine.Use(cors.New())
	timeLine.Use(logger.New())

	routes.ServerRouter(timeLine)

	go func() {
		log.Println("Starting RabbitMQ consumer ....")
		err := consumer.Listen([]string{""})
		if err != nil {
			log.Fatalf("Failed to listen for RabbitMQ message: %v", err)
		}
	}()

	if err := timeLine.Listen(":80"); err != nil {
		log.Fatalf("Unable to listen timeline service %s", err.Error())
	}
}
