package models

// Error is a struct containing an error message.
type Error struct {
	Message string `json:"message"`
}

// String returns the Message from an error.
func (e *Error) String() string {
	return e.Message
}
