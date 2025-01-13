package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gfteix/book_loan_system/pkg/config"
	"github.com/gfteix/book_loan_system/pkg/db"
	"github.com/gfteix/book_loan_system/pkg/mq"
	"github.com/gfteix/book_loan_system/types"
	"github.com/google/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := db.NewPostgreSQLStorage(db.DBConfig{
		DBHost:     config.Envs.DBHost,
		DBPort:     config.Envs.DBPort,
		DBUser:     config.Envs.DBUser,
		DBName:     config.Envs.DBName,
		DBPassword: config.Envs.DBPassword,
	})

	if err != nil {
		log.Fatalf("error starting db: %v", err)
	}

	loans, err := getLoansToProcess(db)

	if err != nil {
		log.Fatalf("error getting loans: %v", err)
	}

	qty := len(loans)

	log.Printf("loan-reminder: processing %v loans", qty)

	if qty > 0 {
		process(loans)
	}
}

func buildPayload(loan types.Loan, eventType string) []byte {
	payload, err := json.Marshal(types.Event{
		Source:  "loan-reminder",
		Time:    time.Now().UTC().Format(time.RFC3339),
		EventId: uuid.NewString(),
		Type:    eventType,
		Payload: types.EventPayload{
			UserId: loan.UserId,
			LoanId: loan.Id,
		},
	})

	if err != nil {
		log.Fatalf("error marshiling json %v", err)
	}

	return payload
}

func publishMessage(ch *amqp.Channel, ctx context.Context, loan types.Loan, queue string) {
	body := buildPayload(loan, queue)

	err := ch.PublishWithContext(ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Fatalf("fail to publish message %v", err)
	}

	log.Printf("[x] Sent %s\n", body)
}

func process(loans []types.Loan) {
	conn, ch, err := mq.NewRabbitMQClient(mq.MQConfig{
		Username: config.Envs.MQUsername,
		Password: config.Envs.MQPassword,
		Host:     config.Envs.MQHost,
		Port:     config.Envs.MQPort,
	})

	if err != nil {
		log.Fatalf("error creating mq client %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	_, err = mq.DeclareQueue(ch, mq.QueueLoanExpired)

	if err != nil {
		log.Fatalf("error declaring queue %v", err)
	}

	_, err = mq.DeclareQueue(ch, mq.QueueLoanExpiring)

	if err != nil {
		log.Fatalf("error declaring queue %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, l := range loans {

		today := time.Now().Truncate(24 * time.Hour)
		expiringDate := l.ExpiringDate.Truncate(24 * time.Hour)

		daysDiff := expiringDate.Sub(today).Hours() / 24

		if daysDiff == 0 {
			publishMessage(ch, ctx, l, mq.QueueLoanExpired)
		}

		if daysDiff > 0 && daysDiff <= 2 {
			publishMessage(ch, ctx, l, mq.QueueLoanExpiring)
		}
	}
}

// returns loans that expires today or that will expire in the next two days

func getLoansToProcess(db *sql.DB) ([]types.Loan, error) {
	rows, err := db.Query("SELECT id, expiring_date, user_id, book_item_id FROM loans WHERE return_date IS NULL AND expiring_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '2 days'")

	if err != nil {
		return nil, err
	}

	loans := make([]types.Loan, 0)

	for rows.Next() {
		loan := new(types.Loan)

		err := rows.Scan(
			&loan.Id,
			&loan.ExpiringDate,
			&loan.UserId,
			&loan.BookItemId,
		)

		if err != nil {
			return nil, err
		}
		loans = append(loans, *loan)
	}

	return loans, nil
}
