package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gfteix/book_loan_system/pkg/config"
	"github.com/gfteix/book_loan_system/pkg/db"

	"github.com/gfteix/book_loan_system/internal/books"
	"github.com/gfteix/book_loan_system/internal/users"
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

	addr := fmt.Sprintf(":%v", config.Envs.Port)
	server := NewAPIServer(addr, db)

	if err := server.Run(); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	userRepository := users.NewRepository(s.db)
	userHandler := users.NewHandler(userRepository)
	userHandler.RegisterRoutes(router)

	bookRepository := books.NewRepository(s.db)
	bookHandler := books.NewHandler(bookRepository)
	bookHandler.RegisterRoutes(router)

	log.Printf("Listening on %v", s.addr)

	return http.ListenAndServe(s.addr, router)
}
