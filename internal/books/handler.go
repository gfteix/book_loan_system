package books

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gfteix/book_loan_system/pkg/utils"
	"github.com/gfteix/book_loan_system/types"
	"github.com/go-playground/validator"
)

type Handler struct {
	repository types.BookRepository
}

func NewHandler(repository types.BookRepository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /books", h.handleCreateBook)
	router.HandleFunc("GET /books", h.handleGetBooks) // add filter category, title, author, isbn, name
	router.HandleFunc("GET /books/{id}", h.handleGetBooksById)
	router.HandleFunc("POST /books/{id}/items", h.handleCreateBookItem)
	router.HandleFunc("GET /books/{id}/items", h.handleGetBookItems) // add filters status, return date, expiring date
	router.HandleFunc("GET /books/{id}/items/{id}", h.handleGetBookItemById)
}

func (h *Handler) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	log.Print("handleGetUserById")

	w.Write([]byte("ola"))

}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Print("handleCreateUser")

	var payload types.CreateUserPayload

	err := utils.ParseJson(r, &payload)

	if err != nil {
		log.Printf("error on ParseJson %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	user, err := h.repository.GetUserByEmail(payload.Email)

	if err != nil {
		log.Printf("error on GetUserByEmail %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if user != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	err = h.repository.CreateUser(types.User{
		Email: payload.Email,
		Name:  payload.Name,
	})

	if err != nil {
		log.Printf("error on CreateUser %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
