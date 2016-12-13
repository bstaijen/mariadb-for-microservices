package models

type Error struct {
	Message string `json:"message"`
}

func (e *Error) String() string {
	return e.Message
}
