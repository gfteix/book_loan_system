package types

import "time"

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

type BookItem struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	BookId    string    `json:"bookId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type Loan struct {
	Id                 string    `json:"id"`
	UserId             string    `json:"userId"`
	BookItemId         string    `json:"bookItemId"`
	Status             string    `json:"status"`
	ExpectedReturnDate time.Time `json:"expectedReturnDate"`
	ReturnDate         time.Time `json:"returnDate"`
	LoanDate           time.Time `json:"loanDate"`
	CreatedAt          time.Time `json:"createdAt"`
}

type UserRepository interface {
	GetUserById(id string) (*User, error)
	GetUserByEmail(id string) (*User, error)
	CreateUser(user User) error
}

type CreateUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
