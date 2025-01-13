package mq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueLoanExpiring string = "loan-expiring"
	QueueLoanExpired  string = "loan-expired"
)

func setQueues(ch *amqp.Channel) {
	_, err := ch.QueueDeclare(
		QueueLoanExpiring, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)

	failOnError(err, fmt.Sprintf("Failed to declare %v queue", QueueLoanExpiring))

	_, err = ch.QueueDeclare(
		QueueLoanExpired, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)

	failOnError(err, fmt.Sprintf("Failed to declare %v queue", QueueLoanExpired))
}

func NewRabbitMQClient() (*amqp.Channel, context.Context) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setQueues(ch)

	return ch, ctx
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
