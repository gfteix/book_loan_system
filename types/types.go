package types

import (
	"context"
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Book struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ISBN          string    `json:"isbn"`
	Author        string    `json:"author"`
	NumberOfPages int       `json:"numberOfPages"`
	CreatedAt     time.Time `json:"createdAt"`
}

type BookCopy struct {
	Id        string    `json:"id"`
	BookId    string    `json:"bookId"`
	Location  string    `json:"location"`
	Condition string    `json:"condition"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type Loan struct {
	Id           string     `json:"id"`
	UserId       string     `json:"userId"`
	BookCopyId   string     `json:"bookCopyId"`
	Status       string     `json:"status"`
	ExpiringDate time.Time  `json:"expiringDate"`
	ReturnDate   *time.Time `json:"returnDate,omitempty"`
	LoanDate     time.Time  `json:"loanDate"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type UserRepository interface {
	GetUserById(id string) (*User, error)
	GetUserByEmail(id string) (*User, error)
	CreateUser(user User) error
}

type BookRepository interface {
	GetBookById(id string) (*Book, error)
	GetBooks(filter map[string]string) ([]Book, error)
	GetBookCopiesByBookId(id string) ([]BookCopy, error)
	GetBookCopyById(id string) (*BookCopy, error)
	CreateBook(book Book) error
	CreateBookCopy(bookCopy BookCopy) error
}

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan Loan) error
	GetLoan(id string) (*Loan, error)
	GetLoans(filters map[string]string) ([]Loan, error)
}

type CreateUserPayload struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type CreateBookPayload struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ISBN          string `json:"isbn"`
	Author        string `json:"author"`
	NumberOfPages int    `json:"numberOfPages"`
}

type CreateBookCopyPayload struct {
	BookId    string `json:"bookId"`
	Status    string `json:"status"`
	Condition string `json:"condition"`
	Location  string `json:"location"`
}

type CreateLoanPayload struct {
	UserId       string    `json:"userId"`
	BookCopyId   string    `json:"bookCopyId"`
	Status       string    `json:"status"`
	ExpiringDate time.Time `json:"expiringDate"`
	LoanDate     time.Time `json:"loanDate"`
}

type EventPayload struct {
	UserId string `json:"userId"`
	LoanId string `json:"loanId"`
}

type Event struct {
	Source  string       `json:"source"`
	Time    string       `json:"time"`
	EventId string       `json:"eventId"`
	Type    string       `json:"type"`
	Payload EventPayload `json:"payload"`
}

type APIError struct {
	Error string `json:"error"`
}
