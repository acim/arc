package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ErrEmptyPassword is returned if user supplied an empty password.
var ErrEmptyPassword = errors.New("empty password")

// User model.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:",omitempty"`
}

// NewUser creates new user model.
func NewUser(email, password string) (*User, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("new uuid: %w", err)
	}

	u := &User{
		ID:       uuid.String(),
		Email:    email,
		Password: password,
	}

	err = u.HashPassword()
	if err != nil {
		return nil, fmt.Errorf("new user: %w", err)
	}

	return u, nil
}

// IsValidPassword returns true if provided plain password matches stored hash.
func (u *User) IsValidPassword(plainPassword string) bool {
	if u.Password == "" || plainPassword == "" {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)) == nil
}

// HashPassword hashes password.
func (u *User) HashPassword() error {
	if u.Password == "" {
		return fmt.Errorf("hash password: %w", ErrEmptyPassword)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	u.Password = string(hash)

	return nil
}
