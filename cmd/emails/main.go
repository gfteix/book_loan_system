package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/gfteix/book_loan_system/pkg/config"
	"github.com/gfteix/book_loan_system/pkg/db"
	"github.com/gfteix/book_loan_system/pkg/mq"
	"github.com/gfteix/book_loan_system/types"
	"github.com/rabbitmq/amqp091-go"
)

type LoanData struct {
	Email         string
	Expiring_date time.Time
	BookTitle     string
}

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

	log.Print("successfully connected to rabbit mq client")

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
	ch.Close() // Stop receiving new messages
	wg.Wait()  // Wait for all in-flight messages to be processed
	log.Println("All workers finished")
}

func processMessage(d amqp091.Delivery) {
	log.Printf("Received message from LoanEVents: %s", d.Body)

	var body types.Event

	err := json.Unmarshal(d.Body, &body)

	if err != nil {
		log.Printf("fail to unmarshal message body %v", err)
	}

	validTypes := []string{"LoanExpired", "LoanExpiring"}

	if !slices.Contains(validTypes, body.Type) {
		log.Printf("Unrecognized event type: %s", body.Type)
		return
	}

	data, err := getDataForEmail(body.Payload.LoanId)

	if err != nil {
		log.Printf("error on getDataForEmail: %v", err)
		return
	}

	var subject string
	var message string

	switch body.Type {
	case "LoanExpired":
		subject = "Loan Expired"
		message = fmt.Sprintf("Your loan of the book %v expired on %v, please return the book to the library.",
			data.BookTitle, data.Expiring_date.Format("2006-01-02"))
	case "LoanExpiring":
		subject = "Loan Expiring"
		message = fmt.Sprintf(
			"Your loan of the book %v will expire on %v, please remember to return the book to the library until the expiration date.",
			data.BookTitle, data.Expiring_date.Format("2006-01-02"))
	}

	err = sendEmail([]string{data.Email}, subject, message)

	if err != nil {
		log.Printf("error processing message %v", err)
		return
	}
	d.Ack(false)
}

func getDataForEmail(loanId string) (LoanData, error) {
	db, err := db.NewPostgreSQLStorage(db.DBConfig{
		DBHost:     config.Envs.DBHost,
		DBPort:     config.Envs.DBPort,
		DBUser:     config.Envs.DBUser,
		DBName:     config.Envs.DBName,
		DBPassword: config.Envs.DBPassword,
	})
	if err != nil {
		return LoanData{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT u.email, l.expiring_date, b.title FROM loans l INNER JOIN users u ON l.user_id = u.id INNER JOIN book_items bi ON bi.id = l.book_item_id INNER JOIN books b ON b.id = bi.book_id WHERE l.id = $1", loanId)

	if err != nil {
		return LoanData{}, err
	}

	defer rows.Close()

	var data LoanData

	if rows.Next() {
		err := rows.Scan(
			&data.Email,
			&data.Expiring_date,
			&data.BookTitle,
		)
		if err != nil {
			return data, err
		}
	}

	return data, nil
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
