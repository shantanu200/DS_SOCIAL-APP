package notifications

import (
	"github.com/gofiber/websocket/v2"
)

type LikeNotificationObj struct {
	AuthorId int `json:"authorId"`
	SenderId int `json:"senderId"`
}

func SendNotification(conn *websocket.Conn, notification string) error {
	if conn != nil {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(notification)); err != nil {
			return err
		}
	}

	return nil
}
