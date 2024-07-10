package rabbitmq

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn     *amqp.Connection
	ConsumeChannel *amqp.Channel
	EXCHANGENAME   = "fanout_exchange"
	QueueName      string
)

func Connect() {
	var tries int64
	backOff := 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

		if err != nil {
			fmt.Println("RabbitMQ is not ready...")
			tries++
		} else {
			log.Println("RabbitMQ connected!!")
			connection = c
			break
		}

		if tries > 5 {
			fmt.Println(err)
			panic(err)
		}

		backOff = time.Duration(math.Pow(float64(tries), 2)) * time.Second
		log.Println("Backing OFF!!")
		time.Sleep(backOff)
		continue
	}

	RabbitConn = connection
}

func SetupConsumeQueue() {
	ch, err := RabbitConn.Channel()
	if err != nil {
		panic(err)
	}


	err = ch.ExchangeDeclare(EXCHANGENAME, "fanout", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)

	QueueName = q.Name

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(q.Name, "", EXCHANGENAME, false, nil)
	if err != nil {
		panic(err)
	}

	ConsumeChannel = ch
}
