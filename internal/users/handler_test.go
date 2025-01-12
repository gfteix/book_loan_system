package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gfteix/book_loan_system/types"
)

type mockUserRepository struct {
	GetUserByEmailFunc func(email string) (*types.User, error)
	GetUserByIdFunc    func(id string) (*types.User, error)
	CreateUserFunc     func(user types.User) error
}

func TestCreateUserHandler(t *testing.T) {
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

	t.Run("should return 400 if payload is missing required fields", func(t *testing.T) {
		userRepository := &mockUserRepository{}

		handler := NewHandler(userRepository)

		payload := map[string]string{
			"Email": "newuser@example.com", // Missing "Name"
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

	t.Run("should fail if email is already registered", func(t *testing.T) {
		userRepository := &mockUserRepository{
			GetUserByEmailFunc: func(email string) (*types.User, error) {
				return &types.User{Id: "1", Email: "existing@email.com"}, nil
			},
		}

		handler := NewHandler(userRepository)

		payload := types.CreateUserPayload{
			Name:  "user",
			Email: "existing@email.com",
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

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return 500 if repository fails on CreateUser", func(t *testing.T) {
		userRepository := &mockUserRepository{
			GetUserByEmailFunc: func(email string) (*types.User, error) {
				return nil, nil // No user exists with this email
			},
			CreateUserFunc: func(user types.User) error {
				return fmt.Errorf("database error")
			},
		}

		handler := NewHandler(userRepository)

		payload := types.CreateUserPayload{
			Name:  "user",
			Email: "newuser@example.com",
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

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

}

func TestGetUserByIdHandler(t *testing.T) {
	t.Run("should fail if user ID is invalid", func(t *testing.T) {
		userRepository := &mockUserRepository{}

		handler := NewHandler(userRepository)

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/users/{id}", handler.handleGetUserById)

		req, err := http.NewRequest(http.MethodGet, "/users/invalid-id", nil)

		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return 404 if user ID does not exist", func(t *testing.T) {
		userRepository := &mockUserRepository{
			GetUserByIdFunc: func(id string) (*types.User, error) {
				return nil, nil // Simulate user not found
			},
		}

		handler := NewHandler(userRepository)

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/users/{id}", handler.handleGetUserById)

		req, err := http.NewRequest(http.MethodGet, "/users/123e4567-e89b-12d3-a456-426614174000", nil)

		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should return 500 if repository fails on GetUserById", func(t *testing.T) {
		userRepository := &mockUserRepository{
			GetUserByIdFunc: func(id string) (*types.User, error) {
				return nil, fmt.Errorf("database error")
			},
		}

		handler := NewHandler(userRepository)

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/users/{id}", handler.handleGetUserById)

		req, err := http.NewRequest(http.MethodGet, "/users/123e4567-e89b-12d3-a456-426614174000", nil)

		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should retrieve user successfully", func(t *testing.T) {
		userRepository := &mockUserRepository{
			GetUserByIdFunc: func(id string) (*types.User, error) {
				return &types.User{
					Id:    id,
					Name:  "Test User",
					Email: "test@example.com",
				}, nil
			},
		}

		handler := NewHandler(userRepository)

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/users/{id}", handler.handleGetUserById)

		req, err := http.NewRequest(http.MethodGet, "/users/123e4567-e89b-12d3-a456-426614174000", nil)

		if err != nil {
			t.Fatal(err)
		}

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func (m *mockUserRepository) GetUserByEmail(email string) (*types.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(email)
	}
	return nil, nil
}

func (m *mockUserRepository) GetUserById(id string) (*types.User, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(id)
	}
	return nil, nil
}

func (m *mockUserRepository) CreateUser(user types.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return nil
}
