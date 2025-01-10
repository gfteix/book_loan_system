package types

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"firstName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserRepository interface {
	GetUserById(id int) (*User, error)
}
