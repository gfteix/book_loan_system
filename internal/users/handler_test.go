package users

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gfteix/book_loan_system/types"
)

type mockUserRepository struct{}

func TestUserServiceHandlers(t *testing.T) {

	userRepository := &mockUserRepository{}
	handler := NewHandler(userRepository)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.CreateUserPayload{
			Name:  "user",
			Email: "invalid",
		}

		marshalled, _ := json.Marshal(payload)

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/users", handler.handleCreateUser)

		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.CreateUserPayload{
			Name:  "user",
			Email: "valid@email.com",
		}

		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/users", handler.handleCreateUser)

		router.ServeHTTP(rr, req)

		log.Print(rr.Body.String())
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

func (m *mockUserRepository) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil // returning no user
}

func (m *mockUserRepository) GetUserById(id string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserRepository) CreateUser(user types.User) error {
	return nil
}
