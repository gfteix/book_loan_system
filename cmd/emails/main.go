package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gfteix/book_loan_system/pkg/config"
	"github.com/gfteix/book_loan_system/pkg/db"
	"github.com/gfteix/book_loan_system/pkg/mq"
	"github.com/gfteix/book_loan_system/types"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
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

	var wg sync.WaitGroup

	go func() {
		for d := range messages {
			wg.Add(1)
			go func(d amqp091.Delivery) {
				defer wg.Done()
				processMessage(d)
			}(d)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	log.Println("Shutting down gracefully...")
	ch.Close() // Close RabbitMQ channel
	wg.Wait()  // Wait for all message handlers to finish
	log.Println("All workers finished")
}

func processMessage(d amqp091.Delivery) {
	log.Printf("Received message from LoanEVents: %s", d.Body)

	var body types.Event

	err := json.Unmarshal(d.Body, &body)

	if err != nil {
		log.Printf("fail to unmarshal message body %v", err)
	}

	switch body.Type {
	case "LoanExpired":
		err = handleExpiredLoan(body)
	case "LoanExpiring":
		err = handleExpiringLoan(body)
	default:
		log.Printf("Unrecognized event type: %s", body.Type)
		return
	}

	if err != nil {
		log.Printf("error processing message %v", err)
		return
	}
	d.Ack(false)
}

func getDataForEmail(loanId string) (string, time.Time, error) {
	db, err := db.NewPostgreSQLStorage(db.DBConfig{
		DBHost:     config.Envs.DBHost,
		DBPort:     config.Envs.DBPort,
		DBUser:     config.Envs.DBUser,
		DBName:     config.Envs.DBName,
		DBPassword: config.Envs.DBPassword,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT u.email, l.expiring_date FROM loans l INNER JOIN users u ON l.user_id = u.id WHERE l.id = $1", loanId)

	if err != nil {
		return "", time.Time{}, err
	}

	defer rows.Close()

	var email string
	var expiring_date time.Time

	if rows.Next() {
		err := rows.Scan(
			&email,
			&expiring_date,
		)
		if err != nil {
			return "", time.Time{}, err
		}
	}

	return email, expiring_date, nil
}

func handleExpiredLoan(event types.Event) error {
	email, expiring_date, err := getDataForEmail(event.Payload.LoanId)

	if err != nil {
		return err
	}

	err = sendEmail([]string{email}, "Loan Expired", fmt.Sprintf("Your book loan expired on %v", expiring_date))

	if err != nil {
		return err
	}
	return nil
}

func handleExpiringLoan(event types.Event) error {
	email, expiring_date, err := getDataForEmail(event.Payload.LoanId)

	if err != nil {
		return err
	}

	err = sendEmail([]string{email}, "Loan Expiring", fmt.Sprintf("Your book loan will expire on %v", expiring_date))

	if err != nil {
		return err
	}
	return nil
}

func sendEmail(to []string, subject string, body string) error {
	from := "book_loan_system@email.com"
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	auth := smtp.PlainAuth("", "", "", config.Envs.SMTPHost)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", config.Envs.SMTPHost, config.Envs.SMTPPort), auth, from, to, []byte(msg))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
