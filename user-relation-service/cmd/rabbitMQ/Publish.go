package rabbitmq

import (
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishPayload struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type FollowerPayload struct {
	IsFollow   bool  `json:"isFollow"`
	FolloweeId int64 `json:"followeeId"`
	FollowerId int64 `json:"followerId"`
}

func PublishQueue(ExchangeName string, Event string, message string) error {

	payload := PublishPayload{
		Event: Event,
		Data:  message,
	}

	jsonMarshal, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	err = ConsumeChannel.Publish(ExchangeName, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        jsonMarshal,
	})

	if err != nil {
		return err
	}

	return err
}
