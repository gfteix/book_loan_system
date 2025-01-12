package books

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gfteix/book_loan_system/types"
)

type mockBookRepository struct {
	GetBookByIdFunc          func(id string) (*types.Book, error)
	GetBooksFunc             func(filter map[string]string) ([]types.Book, error)
	CreateBookFunc           func(book types.Book) error
	CreateBookItemFunc       func(bookItem types.BookItem) error
	GetBookItemsByBookIdFunc func(bookId string) ([]types.BookItem, error)
	GetBookItemByIdFunc      func(itemId string) (*types.BookItem, error)
}

func (m *mockBookRepository) GetBookById(id string) (*types.Book, error) {
	if m.GetBookByIdFunc != nil {
		return m.GetBookByIdFunc(id)
	}
	return nil, nil
}

func (m *mockBookRepository) GetBooks(filter map[string]string) ([]types.Book, error) {
	if m.GetBooksFunc != nil {
		return m.GetBooksFunc(filter)
	}
	return nil, nil
}

func (m *mockBookRepository) CreateBook(book types.Book) error {
	if m.CreateBookFunc != nil {
		return m.CreateBookFunc(book)
	}
	return nil
}

func (m *mockBookRepository) CreateBookItem(bookItem types.BookItem) error {
	if m.CreateBookItemFunc != nil {
		return m.CreateBookItemFunc(bookItem)
	}
	return nil
}

func (m *mockBookRepository) GetBookItemsByBookId(bookId string) ([]types.BookItem, error) {
	if m.GetBookItemsByBookIdFunc != nil {
		return m.GetBookItemsByBookIdFunc(bookId)
	}
	return nil, nil
}

func (m *mockBookRepository) GetBookItemById(itemId string) (*types.BookItem, error) {
	if m.GetBookItemByIdFunc != nil {
		return m.GetBookItemByIdFunc(itemId)
	}
	return nil, nil
}

func TestBookHandler(t *testing.T) {
	repository := &mockBookRepository{}
	handler := NewHandler(repository)

	t.Run("should fail if creating a book with invalid payload", func(t *testing.T) {
		payload := map[string]interface{}{
			"title": 123, // Invalid type for title
		}

		marshalled, _ := json.Marshal(payload)
		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/books", handler.handleCreateBook)

		req, err := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should successfully create a book", func(t *testing.T) {
		repository.CreateBookFunc = func(book types.Book) error {
			return nil
		}

		payload := types.CreateBookPayload{
			Title:         "Sample Book",
			Description:   "A book description",
			ISBN:          "1234567890",
			Author:        "Author Name",
			NumberOfPages: 100,
		}

		marshalled, _ := json.Marshal(payload)
		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/books", handler.handleCreateBook)

		req, err := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should fail to fetch a book if not found", func(t *testing.T) {
		repository.GetBookByIdFunc = func(id string) (*types.Book, error) {
			return nil, nil
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/books/{id}", handler.handleGetBookById)

		req, err := http.NewRequest(http.MethodGet, "/books/123e4567-e89b-12d3-a456-426614174000", nil)
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should fetch all books successfully", func(t *testing.T) {
		repository.GetBooksFunc = func(filter map[string]string) ([]types.Book, error) {
			return []types.Book{
				{Title: "Book 1", Author: "Author 1"},
				{Title: "Book 2", Author: "Author 2"},
			}, nil
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/books", handler.handleGetBooks)

		req, err := http.NewRequest(http.MethodGet, "/books", nil)
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fetch books using filter params", func(t *testing.T) {
		filterKey := "title"
		filterValue := "Book1"

		var gotValue string

		repository.GetBooksFunc = func(filter map[string]string) ([]types.Book, error) {
			gotValue = filter[filterKey]

			return []types.Book{
				{Title: "Book1", Author: "Author 1"},
			}, nil
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		path := fmt.Sprintf("/books?%v=%v", filterKey, filterValue)

		router.HandleFunc("/books", handler.handleGetBooks)

		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if gotValue != filterValue {
			t.Errorf("expected filter value %v, got %v", filterValue, gotValue)
		}
	})

	t.Run("should fail to create book item with invalid book ID", func(t *testing.T) {
		payload := types.CreateBookItem{
			BookId:    "invalid-id",
			Status:    "Available",
			Location:  "Library",
			Condition: "New",
		}

		marshalled, _ := json.Marshal(payload)
		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/books/{id}/items", handler.handleCreateBookItem)

		req, err := http.NewRequest(http.MethodPost, "/books/invalid-id/items", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fetch book items successfully", func(t *testing.T) {
		repository.GetBookItemsByBookIdFunc = func(bookId string) ([]types.BookItem, error) {
			return []types.BookItem{
				{BookId: "book-id", Status: "Available", Location: "Library"},
			}, nil
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/books/{id}/items", handler.handleGetBookItems)

		req, err := http.NewRequest(http.MethodGet, "/books/book-id/items", nil)
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
