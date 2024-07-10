package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"notification/cmd/api/routes"
	"notification/cmd/cache"
	"notification/cmd/database"
	"notification/cmd/socket/notifications"

	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

type Notification struct {
	Message string                 `json:"message"`
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type client struct {
	isClosing bool
	mu        sync.Mutex
}

var (
	clients         = make(map[*websocket.Conn]*client)
	userConnections = make(map[string]*websocket.Conn)
	register        = make(chan *websocket.Conn)
	broadcast       = make(chan Notification)
	unregister      = make(chan *websocket.Conn)
)

var ctx = context.Background()

func runHub() {
	for {
		select {
		case connection := <-register:
			userId := connection.Query("id")
			userConnections[userId] = connection
			clients[connection] = &client{}
			log.Printf("Connection registered with userId %s: ", userId)

		case notification := <-broadcast:
			cacheKey := "notifications"
			if notification.Type == "LIKE" {
				payload := notification.Payload
				authorId := fmt.Sprint(payload["authorId"])
				fieldKey := notification.Type + ":" + authorId + ":" + fmt.Sprint(payload["postId"])
				cacheGet := cache.RedisClient.HGet(ctx, cacheKey, fieldKey)

				if cacheGet.Err() == nil {
					fmt.Println("Already notification send")
				} else {
					wsConnection := userConnections[authorId]

					if wsConnection.Query("id") == authorId {
						fmt.Println("Author liked his/her post")
					}

					notificationJSON, err := json.Marshal(notification)

					if err != nil {
						fmt.Println(err)
					}

					err = notifications.SendNotification(wsConnection, string(notificationJSON))

					cacheSet := cache.RedisClient.HSet(ctx, cacheKey, fieldKey, true)

					if cacheSet.Err() != nil {
						fmt.Println(cacheSet.Err())
					}

					if err != nil {
						wsConnection.WriteMessage(websocket.CloseMessage, []byte{})
						wsConnection.Close()
						unregister <- wsConnection
					}
				}

			}

			if notification.Type == "POST" {
				payload := notification.Payload
				authorId := fmt.Sprint(payload["authorId"])

				notificationJSON, err := json.Marshal(notification)

				if err != nil {
					log.Fatal(err)
				}

				cacheKey := "follower:user:" + authorId
				cacheGet := cache.RedisClient.SMembers(ctx, cacheKey)

				if cacheGet.Err() != nil {
					log.Fatal(cacheGet.Err())
				}

				for _, ids := range cacheGet.Val() {
					go func(connection *websocket.Conn, c *client) {
						c.mu.Lock()
						defer c.mu.Unlock()

						if c.isClosing {
							return
						}

						err = notifications.SendNotification(connection, string(notificationJSON))

						if err != nil {
							fmt.Println(err)
						}

					}(userConnections[ids], clients[userConnections[ids]])
				}
			}

			if notification.Type == "FOLLOWER" {
				payload := notification.Payload
				followeeId := fmt.Sprint(payload["followerId"])

				connection := userConnections[followeeId]

				notificationJSON, err := json.Marshal(notification)

				if err != nil {
					fmt.Println(err)
				}

				err = notifications.SendNotification(connection, string(notificationJSON))

				if err != nil {
					fmt.Println(err)
				}
			}

		case connection := <-unregister:
			delete(clients, connection)
			log.Println("Connection unregistered")
		}
	}
}

func main() {
	app := fiber.New()
	app.Use(healthcheck.New())
	app.Use(cors.New())
	app.Use(logger.New())

	database.ConnectPostgres()
	cache.ConnectCache()

	routes.ServerRouter(app)

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go runHub()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			c.Close()
		}()

		register <- c

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			var notification Notification
			if err := json.Unmarshal(msg, &notification); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}

			broadcast <- notification
		}
	}))

	addr := flag.String("addr", ":80", "HTTP service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
