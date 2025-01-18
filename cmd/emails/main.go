package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gfteix/book_loan_system/pkg/config"
	"github.com/gfteix/book_loan_system/pkg/db"
	"github.com/gfteix/book_loan_system/pkg/mq"
)

func main() {
	_, err := db.NewPostgreSQLStorage(db.DBConfig{
		DBHost:     config.Envs.DBHost,
		DBPort:     config.Envs.DBPort,
		DBUser:     config.Envs.DBUser,
		DBName:     config.Envs.DBName,
		DBPassword: config.Envs.DBPassword,
	})

	if err != nil {
		log.Fatalf("error starting db: %v", err)
	}

	conn, ch, err := mq.NewRabbitMQClient(mq.MQConfig{
		Username: config.Envs.MQUsername,
		Password: config.Envs.MQPassword,
		Host:     config.Envs.MQHost,
		Port:     config.Envs.MQPort,
	})

	if err != nil {
		log.Fatalf("error creating mq client: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	_, err = mq.DeclareQueue(ch, "LoanEvents")
	if err != nil {
		log.Fatalf("error declaring queue: %v", err)
	}

	messages, err := ch.Consume("LoanEvents", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error consuming LoanEvents: %v", err)
	}

	go func() {
		for d := range messages {
			log.Printf("Received message from LoanEVents: %s", d.Body)
			d.Ack(false)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-stopChan
	log.Println("Shutting down gracefully...")
}
