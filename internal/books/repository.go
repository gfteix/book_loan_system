package books

import (
	"database/sql"

	"github.com/gfteix/book_loan_system/types"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func scanRowIntoBook(rows *sql.Rows) (*types.Book, error) {
	book := new(types.Book)
	err := rows.Scan(
		&book.Id,
		&book.Title,
		&book.Description,
		&book.ISBN,
		&book.Author,
		&book.NumberOfPages,
		&book.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func scanRowIntoBookItem(rows *sql.Rows) (*types.BookItem, error) {
	bookItem := new(types.BookItem)
	err := rows.Scan(
		&bookItem.Id,
		&bookItem.UserId,
		&bookItem.BookId,
		&bookItem.Status,
		&bookItem.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return bookItem, nil
}

func scanRowIntoLoan(rows *sql.Rows) (*types.Loan, error) {
	loan := new(types.Loan)
	err := rows.Scan(
		&loan.Id,
		&loan.UserId,
		&loan.BookItemId,
		&loan.Status,
		&loan.ExpiringDate,
		&loan.ReturnDate,
		&loan.LoanDate,
		&loan.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (r *Repository) GetBookById(id string) (*types.Book, error) {
	rows, err := r.db.Query("SELECT id, title, description, isbn, author, number_of_pages, created_at FROM books WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanRowIntoBook(rows)
	}

	return nil, nil
}

func (r *Repository) GetBookItemsByBookId(id string) ([]types.BookItem, error) {
	rows, err := r.db.Query("SELECT id, user_id, book_id, status, created_at FROM book_items WHERE book_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookItems []types.BookItem
	for rows.Next() {
		bookItem, err := scanRowIntoBookItem(rows)
		if err != nil {
			return nil, err
		}
		bookItems = append(bookItems, *bookItem)
	}

	return bookItems, nil
}

func (r *Repository) GetBookItemById(id string) (*types.BookItem, error) {
	rows, err := r.db.Query("SELECT id, user_id, book_id, status, created_at FROM book_items WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanRowIntoBookItem(rows)
	}

	return nil, nil
}

func (r *Repository) CreateBook(book types.Book) error {
	id := uuid.NewString()
	_, err := r.db.Exec("INSERT INTO books (id, title, description, isbn, author, number_of_pages, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, book.Title, book.Description, book.ISBN, book.Author, book.NumberOfPages, book.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateBookItem(bookItem types.BookItem) error {
	id := uuid.NewString()
	_, err := r.db.Exec("INSERT INTO book_items (id, user_id, book_id, status, created_at) VALUES ($1, $2, $3, $4, $5)",
		id, bookItem.UserId, bookItem.BookId, bookItem.Status, bookItem.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateLoan(loan types.Loan) error {
	id := uuid.NewString()
	_, err := r.db.Exec("INSERT INTO loans (id, user_id, book_item_id, status, expiring_date, return_date, loan_date, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, loan.UserId, loan.BookItemId, loan.Status, loan.ExpiringDate, loan.ReturnDate, loan.LoanDate, loan.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLoan(id string) (*types.Loan, error) {
	rows, err := r.db.Query("SELECT id, user_id, book_item_id, status, expiring_date, return_date, loan_date, created_at FROM loans WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loan *types.Loan

	for rows.Next() {
		loan, err = scanRowIntoLoan(rows)
		if err != nil {
			return nil, err
		}
	}

	return loan, nil
}

func (r *Repository) GetLoans() ([]types.Loan, error) {
	rows, err := r.db.Query("SELECT id, user_id, book_item_id, status, expiring_date, return_date, loan_date, created_at FROM loans")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []types.Loan
	for rows.Next() {
		loan, err := scanRowIntoLoan(rows)
		if err != nil {
			return nil, err
		}
		loans = append(loans, *loan)
	}

	return loans, nil
}
