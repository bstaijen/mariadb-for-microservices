package models

import (
	"errors"
	"fmt"
	"time"
)

// lower_case private, upper_case public
// Uppercase variable is mandatory for exposing to json

// User model
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	Password  string    `json:"password,omitempty"`
	Hash      string
}

// GetUsername returns the username of an user
func (u *User) GetUsername() string {
	return u.Username
}

// Print method of User
func (u *User) Print() string {
	return fmt.Sprintf("%v (%v) - %v", u.Username, u.ID, u.CreatedAt)
}

// Validate returns an error when the username or email is to short.
func (u *User) Validate() error {
	if len(u.Username) < 1 {
		return ErrUsernameTooShort
	}

	if len(u.Email) < 1 {
		return ErrEmailTooShort
	}

	return nil
}

// ValidatePassword returns an error if the password is to small.
func (u *User) ValidatePassword() error {
	if len(u.Password) < 1 {
		return ErrPasswordTooShort
	}

	return nil
}

// ErrUsernameTooShort is a error and is used when username is too short.
var ErrUsernameTooShort = errors.New("Username is too short")

// ErrEmailTooShort is an error and is used when email address is too short.
var ErrEmailTooShort = errors.New("Email address is to short")

// ErrPasswordTooShort is an error and is used when a password is too short.
var ErrPasswordTooShort = errors.New("Password is to short")
