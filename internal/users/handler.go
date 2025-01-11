package users

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
	repository types.UserRepository
}

func NewHandler(repository types.UserRepository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /users", h.handleCreateUser)
	router.HandleFunc("GET /users/{id}", h.handleGetUserById)
}

func (h *Handler) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	log.Print("handleGetUserById")

	id := r.PathValue("id")

	err := uuid.Validate(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	user, err := h.repository.GetUserById(id)

	if err != nil {
		log.Printf("error on handleGetUserById %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
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
