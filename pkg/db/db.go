package db

import (
	"database/sql"

	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"

	"fmt"
	"log"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func NewPostgreSQLStorage(config DBConfig) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := sql.Open("pgx", dataSourceName)

	if err != nil {
		log.Printf("unable to connect to database: %v\n", err)
		return nil, err
	}

	initStorage(db)

	return db, nil
}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Successfuly connected")
}
