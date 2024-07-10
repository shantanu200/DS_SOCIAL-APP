package rabbitmq

import (
	"log"

	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishPayload struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type UpdateUserPayload struct {
	UserId int64 `json:"userId"`
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

	log.Println("EVENT Published to Queue")

	return err
}
