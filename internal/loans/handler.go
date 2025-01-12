package loans

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gfteix/book_loan_system/pkg/utils"
	"github.com/gfteix/book_loan_system/types"
	"github.com/google/uuid"
)

type Handler struct {
	repository types.LoanRepository
}

func NewHandler(repository types.LoanRepository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /loans", h.handleCreateLoan)
	router.HandleFunc("GET /loans", h.handleGetLoans)
	router.HandleFunc("GET /loans/{id}", h.handleGetLoanById)
}

func (h *Handler) handleCreateLoan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload types.CreateLoanPayload

	err := utils.ParseJson(r, &payload)

	if err != nil {
		log.Printf("error on ParseJson %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.repository.CreateLoan(ctx, types.Loan{
		UserId:       payload.UserId,
		BookItemId:   payload.BookItemId,
		Status:       payload.Status,
		ExpiringDate: payload.ExpiringDate,
		LoanDate:     payload.LoanDate,
	})

	if err != nil {
		log.Printf("error on CreateLoan %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetLoans(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	filter := make(map[string]string)

	filter["userId"] = queryParams.Get("userId")
	filter["status"] = queryParams.Get("status")
	filter["bookItemId"] = queryParams.Get("bookItemIds")

	books, err := h.repository.GetLoans(filter)

	if err != nil {
		log.Printf("error on GetBooks %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, books)
}

func (h *Handler) handleGetLoanById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := uuid.Validate(id)

	if err != nil {
		log.Printf("invalid id %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("invalid id"))
		return
	}

	loan, err := h.repository.GetLoan(id)

	if err != nil {
		log.Printf("error on GetLoan %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if loan == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, loan)
}
