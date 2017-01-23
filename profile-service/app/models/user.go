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

// Method of user. (not a function)
// p *user is the receiver of method getUsername()
func (u *User) GetUsername() string {
	return u.Username
}

// Print method of User
func (u *User) Print() string {
	return fmt.Sprintf("%v (%v) - %v", u.Username, u.ID, u.CreatedAt)
}

// Validate method
func (u *User) Validate() error {
	if len(u.Username) < 1 {
		return ErrUsernameTooShort
	}

	if len(u.Email) < 1 {
		return ErrEmailTooShort
	}

	return nil
}

func (u *User) ValidatePassword() error {
	if len(u.Password) < 1 {
		return ErrPasswordTooShort
	}

	return nil
}

var ErrUsernameTooShort = errors.New("Username is too short")
var ErrEmailTooShort = errors.New("Email address is to short")
var ErrPasswordTooShort = errors.New("Password is to short")
