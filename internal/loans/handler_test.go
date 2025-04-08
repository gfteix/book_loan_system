package loans

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gfteix/book_loan_system/types"
)

type mockLoanRepository struct {
	CreateLoanFunc func(ctx context.Context, loan types.Loan) error
	GetLoansFunc   func(filter map[string]string) ([]types.Loan, error)
	GetLoanFunc    func(id string) (*types.Loan, error)
}

func (m *mockLoanRepository) CreateLoan(ctx context.Context, loan types.Loan) error {
	if m.CreateLoanFunc != nil {
		return m.CreateLoanFunc(nil, loan)
	}
	return nil
}

func (m *mockLoanRepository) GetLoans(filter map[string]string) ([]types.Loan, error) {
	if m.GetLoansFunc != nil {
		return m.GetLoansFunc(filter)
	}
	return nil, nil
}

func (m *mockLoanRepository) GetLoan(id string) (*types.Loan, error) {
	if m.GetLoanFunc != nil {
		return m.GetLoanFunc(id)
	}
	return nil, nil
}

func TestLoanHandler(t *testing.T) {
	repository := &mockLoanRepository{}
	handler := NewHandler(repository)

	t.Run("should fail if creating a loan with invalid payload", func(t *testing.T) {
		payload := map[string]interface{}{
			"userId": 123, // Invalid type for userId
		}

		marshalled, _ := json.Marshal(payload)
		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/loans", handler.handleCreateLoan)

		req, err := http.NewRequest(http.MethodPost, "/loans", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should successfully create a loan", func(t *testing.T) {
		repository.CreateLoanFunc = func(ctx context.Context, loan types.Loan) error {
			return nil
		}

		expiringDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
		loanDate := time.Date(2025, time.January, 10, 0, 0, 0, 0, time.UTC)

		payload := types.CreateLoanPayload{
			UserId:       "user-123",
			BookCopyId:   "item-456",
			Status:       "Borrowed",
			ExpiringDate: expiringDate,
			LoanDate:     loanDate,
		}

		marshalled, _ := json.Marshal(payload)
		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/loans", handler.handleCreateLoan)

		req, err := http.NewRequest(http.MethodPost, "/loans", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should fail to fetch a loan if not found", func(t *testing.T) {
		repository.GetLoanFunc = func(id string) (*types.Loan, error) {
			return nil, nil
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/loans/{id}", handler.handleGetLoanById)

		req, err := http.NewRequest(http.MethodGet, "/loans/123e4567-e89b-12d3-a456-426614174000", nil)
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should fetch all loans successfully", func(t *testing.T) {
		repository.GetLoansFunc = func(filter map[string]string) ([]types.Loan, error) {
			return []types.Loan{
				{UserId: "user-1", BookCopyId: "item-1", Status: "Borrowed"},
				{UserId: "user-2", BookCopyId: "item-2", Status: "Returned"},
			}, nil
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/loans", handler.handleGetLoans)

		req, err := http.NewRequest(http.MethodGet, "/loans", nil)
		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fetch loans using filter params", func(t *testing.T) {
		filterKey := "userId"
		filterValue := "user-1"

		var gotValue string

		repository.GetLoansFunc = func(filter map[string]string) ([]types.Loan, error) {
			gotValue = filter[filterKey]
			return []types.Loan{
				{UserId: "user-1", BookCopyId: "item-1", Status: "Borrowed"},
			}, nil
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		path := fmt.Sprintf("/loans?%v=%v", filterKey, filterValue)
		router.HandleFunc("/loans", handler.handleGetLoans)

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
}
