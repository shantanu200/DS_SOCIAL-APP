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

type Consumer struct {
   conn *amqp.Connection
   queueName string
}

func Connect() (*amqp.Connection,error) {
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

    return connection,nil
}

func NewConsumer(conn *amqp.Connection) (Consumer,error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.Setup()

	if err != nil {
		return Consumer{},nil
	}

	return consumer,nil
}

func (consumer *Consumer) Setup() error {
	channel,err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return DeclareExchange(channel)
}




