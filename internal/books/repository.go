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
	_, err := r.db.Exec("INSERT INTO books (id, title, description, isbn, author, number_of_pages) VALUES ($1, $2, $3, $4, $5, $6)",
		id, book.Title, book.Description, book.ISBN, book.Author, book.NumberOfPages)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateBookItem(bookItem types.BookItem) error {
	id := uuid.NewString()
	_, err := r.db.Exec("INSERT INTO book_items (id, user_id, book_id, status) VALUES ($1, $2, $3, $4)",
		id, bookItem.UserId, bookItem.BookId, bookItem.Status)
	if err != nil {
		return err
	}

	return nil
}
