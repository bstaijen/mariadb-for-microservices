package models

import "time"

// User contains user properties
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	Password  string    `json:"password,omitempty"`
}
