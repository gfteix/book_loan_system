package mq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueLoanExpiring string = "loan-expiring"
	QueueLoanExpired  string = "loan-expired"
)

func DeclareQueue(ch *amqp.Channel, name string) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Printf("failed to declare %v queue", QueueLoanExpiring)
		return nil, err
	}

	return &q, nil
}

type MQConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func NewRabbitMQClient(config MQConfig) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/", config.Username, config.Password, config.Host, config.Port))

	if err != nil {
		log.Printf("failed to open rabbit mq connection %v", err)
		return nil, nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Printf("failed to open rabbit mq channel %v", err)
		return nil, nil, err
	}

	return conn, ch, nil
}
