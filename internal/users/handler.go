package users

import (
	"log"
	"net/http"

	"github.com/gfteix/book_loan_system/types"
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

	w.Write([]byte("ola"))

}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Print("handleCreateUser")
}
