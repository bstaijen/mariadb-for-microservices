package models

// CommentCreate is the model for creating a comment
type CommentCreate struct {
	UserID  int    `json:"user_id"`
	PhotoID int    `json:"photo_id"`
	Comment string `json:"comment"`
}
