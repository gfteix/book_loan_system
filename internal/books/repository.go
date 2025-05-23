package books

import (
	"database/sql"
	"fmt"
	"strings"

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

func scanRowIntoBookCopy(rows *sql.Rows) (*types.BookCopy, error) {
	bookCopy := new(types.BookCopy)
	err := rows.Scan(
		&bookCopy.Id,
		&bookCopy.BookId,
		&bookCopy.Status,
		&bookCopy.Location,
		&bookCopy.Condition,
		&bookCopy.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return bookCopy, nil
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

func (r *Repository) GetBooks(filters map[string]string) ([]types.Book, error) {
	q := "SELECT id, title, description, isbn, author, number_of_pages, created_at FROM books"

	where := make([]string, 0)
	whereValues := make([]string, 0)

	whereIndex := 1
	for k, v := range filters {
		if v == "" {
			continue
		}

		if k == "title" {
			where = append(where, fmt.Sprintf("title = $%v", whereIndex))
			whereValues = append(whereValues, v)
			whereIndex++
		}

		if k == "isbn" {
			where = append(where, fmt.Sprintf("isbn = $%v", whereIndex))
			whereValues = append(whereValues, v)
			whereIndex++
		}

		if k == "author" {
			where = append(where, fmt.Sprintf("author = $%v", whereIndex))
			whereValues = append(whereValues, v)
			whereIndex++
		}
	}

	if len(where) > 0 {
		q = fmt.Sprintf("%v WHERE %v", q, strings.Join(where, " AND "))
		q = strings.TrimSuffix(q, " AND ")
	}
	args := make([]interface{}, len(whereValues))

	for i, v := range whereValues {
		args[i] = v
	}

	rows, err := r.db.Query(q, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := make([]types.Book, 0)

	for rows.Next() {
		book, err := scanRowIntoBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, *book)
	}

	return books, nil
}

func (r *Repository) GetBookCopiesByBookId(id string) ([]types.BookCopy, error) {
	rows, err := r.db.Query("SELECT id, book_id, status, location, condition, created_at FROM book_copies WHERE book_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookCopies []types.BookCopy
	for rows.Next() {
		bookCopy, err := scanRowIntoBookCopy(rows)
		if err != nil {
			return nil, err
		}
		bookCopies = append(bookCopies, *bookCopy)
	}

	return bookCopies, nil
}

func (r *Repository) GetBookCopyById(id string) (*types.BookCopy, error) {
	rows, err := r.db.Query("SELECT id, book_id, status, location, condition, created_at FROM book_copies WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanRowIntoBookCopy(rows)
	}

	return nil, nil
}

func (r *Repository) CreateBook(book types.Book) error {
	id := uuid.NewString()
	_, err := r.db.Exec("INSERT INTO books (id, title, description, isbn, author, number_of_pages) VALUES ($1, $2, $3, $4, $5, $6)",
		id, book.Title, book.Description, book.ISBN, book.Author, book.NumberOfPages)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateBookCopy(bookCopy types.BookCopy) error {
	id := uuid.NewString()
	_, err := r.db.Exec("INSERT INTO book_copies (id, book_id, status, location, condition) VALUES ($1, $2, $3, $4, $5)",
		id, bookCopy.BookId, bookCopy.Status, bookCopy.Location, bookCopy.Condition)
	if err != nil {
		return err
	}

	return nil
}
