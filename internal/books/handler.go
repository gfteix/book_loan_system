package books

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gfteix/book_loan_system/pkg/utils"
	"github.com/gfteix/book_loan_system/types"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Handler struct {
	repository types.BookRepository
}

func NewHandler(repository types.BookRepository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /books", h.handleCreateBook)
	router.HandleFunc("GET /books", h.handleGetBooks)
	router.HandleFunc("GET /books/{id}", h.handleGetBookById)
	router.HandleFunc("POST /books/{id}/items", h.handleCreateBookItem)
	router.HandleFunc("GET /books/{id}/items", h.handleGetBookItems)

}

// handleGetBookById godoc
// @Summary Get a book by ID
// @Description Retrieves a book by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} types.Book
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /books/{id} [get]
func (h *Handler) handleGetBookById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	book, err := h.repository.GetBookById(id)
	if err != nil {
		log.Printf("error on GetBookById %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if book == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("book with id %s not found", id))
		return
	}
	utils.WriteJSON(w, http.StatusOK, book)
}

// handleGetBooks godoc
// @Summary Get books with filters
// @Description Retrieves a list of books with optional filters
// @Tags books
// @Accept  json
// @Produce  json
// @Param title query string false "Filter by title"
// @Param author query string false "Filter by author"
// @Param isbn query string false "Filter by ISBN"
// @Success 200 {array} types.Book
// @Failure 500 {object} types.APIError
// @Router /books [get]
func (h *Handler) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filter := map[string]string{
		"title":  queryParams.Get("title"),
		"author": queryParams.Get("author"),
		"isbn":   queryParams.Get("isbn"),
	}

	books, err := h.repository.GetBooks(filter)
	if err != nil {
		log.Printf("error on GetBooks %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, books)
}

// handleCreateBook godoc
// @Summary Create a new book
// @Description Adds a new book to the library system
// @Tags books
// @Accept  json
// @Produce  json
// @Param book body types.CreateBookPayload true "Book details"
// @Success 201
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /books [post]
func (h *Handler) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	log.Print("handleCreateBook")
	var payload types.CreateBookPayload

	err := utils.ParseJson(r, &payload)
	if err != nil {
		log.Printf("error on ParseJson %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err = h.repository.CreateBook(types.Book{
		Title:         payload.Title,
		Description:   payload.Description,
		ISBN:          payload.ISBN,
		Author:        payload.Author,
		NumberOfPages: payload.NumberOfPages,
	})

	if err != nil {
		log.Printf("error on CreateBook %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

// handleCreateBookItem godoc
// @Summary Create a book item
// @Description Adds a new book item to a book
// @Tags books
// @Accept  json
// @Produce  json
// @Param bookItem body types.CreateBookItemPayload true "Book item details"
// @Param id path string true "Book ID"
// @Success 201
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /books/{id}/items [post]
func (h *Handler) handleCreateBookItem(w http.ResponseWriter, r *http.Request) {
	log.Print("handleCreateBookItem")
	var payload types.CreateBookItemPayload

	err := utils.ParseJson(r, &payload)
	if err != nil {
		log.Printf("error on ParseJson %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err = uuid.Validate(payload.BookId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid bookId"))
		return
	}

	book, err := h.repository.GetBookById(payload.BookId)
	if err != nil {
		log.Printf("error on GetBookById %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if book == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid bookId"))
		return
	}

	err = h.repository.CreateBookItem(types.BookItem{
		BookId:    payload.BookId,
		Status:    payload.Status,
		Location:  payload.Location,
		Condition: payload.Condition,
	})

	if err != nil {
		log.Printf("error on CreateBookItem %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

// handleGetBookItems godoc
// @Summary Get items of a book
// @Description Retrieves all items belonging to a book by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {array} types.BookItem
// @Failure 500 {object} types.APIError
// @Router /books/{id}/items [get]
func (h *Handler) handleGetBookItems(w http.ResponseWriter, r *http.Request) {
	bookId := r.PathValue("id")

	bookItems, err := h.repository.GetBookItemsByBookId(bookId)
	if err != nil {
		log.Printf("error on GetBookItemsByBookId %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, bookItems)
}

// handleGetBookItemById godoc
// @Summary Get a book item by ID
// @Description Retrieves a specific book item by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param bookId path string true "Book ID"
// @Param itemId path string true "Book Item ID"
// @Success 200 {object} types.BookItem
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /books/{bookId}/items/{itemId} [get]
func (h *Handler) handleGetBookItemById(w http.ResponseWriter, r *http.Request) {
	itemId := r.PathValue("itemId")

	bookItem, err := h.repository.GetBookItemById(itemId)
	if err != nil {
		log.Printf("error on GetBookItemById %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if bookItem == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("book item with id %s not found", itemId))
		return
	}
	utils.WriteJSON(w, http.StatusOK, bookItem)
}
