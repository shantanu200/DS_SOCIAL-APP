package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		EXCHANGENAME,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}

func DeclareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare("", false, false, true, false, nil)
}
