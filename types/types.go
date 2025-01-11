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
	BookId    string    `json:"bookId"`
	Location  string    `json:"location"`
	Condition string    `json:"condition"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type Loan struct {
	Id           string    `json:"id"`
	UserId       string    `json:"userId"`
	BookItemId   string    `json:"bookItemId"`
	Status       string    `json:"status"`
	ExpiringDate time.Time `json:"expiringDate"`
	ReturnDate   time.Time `json:"returnDate"`
	LoanDate     time.Time `json:"loanDate"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserRepository interface {
	GetUserById(id string) (*User, error)
	GetUserByEmail(id string) (*User, error)
	CreateUser(user User) error
}

type BookRepository interface {
	GetBookById(id string) (*Book, error)
	GetBooks(filter map[string]string) ([]Book, error)
	GetBookItemsByBookId(id string) ([]BookItem, error)
	GetBookItemById(id string) (*BookItem, error)
	CreateBook(book Book) error
	CreateBookItem(bookItem BookItem) error
}

type LoanRepository interface {
	CreateLoan(loan Loan) error
	GetLoan(id string) ([]Loan, error)
	GetLoans() ([]Loan, error)
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

type CreateBookItem struct {
	BookId    string `json:"bookId"`
	Status    string `json:"status"`
	Condition string `json:"condition"`
	Location  string `json:"location"`
}
