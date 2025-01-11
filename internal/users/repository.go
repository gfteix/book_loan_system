package users

import (
	"database/sql"

	"github.com/gfteix/book_loan_system/types"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) CreateUser(user types.User) error {
	id := uuid.NewString()

	_, err := r.db.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", id, user.Name, user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserById(id string) (*types.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, created_at FROM users WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	found := false
	u := new(types.User)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}

		found = true
	}

	if !found {
		return nil, nil
	}

	return u, nil
}

func (r *Repository) GetUserByEmail(email string) (*types.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, created_at FROM users WHERE email = $1", email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	found := false
	u := new(types.User)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}

		found = true
	}

	if !found {
		return nil, nil
	}

	return u, nil
}
