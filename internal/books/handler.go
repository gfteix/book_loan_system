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
	router.HandleFunc("GET /books/{bookId}/items/{itemId}", h.handleGetBookItemById)
}

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

func (h *Handler) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	filter := make(map[string]string)

	filter["title"] = queryParams.Get("title")
	filter["author"] = queryParams.Get("author")
	filter["isbn"] = queryParams.Get("isbn")

	books, err := h.repository.GetBooks(filter)

	if err != nil {
		log.Printf("error on GetBooks %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, books)
}

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
