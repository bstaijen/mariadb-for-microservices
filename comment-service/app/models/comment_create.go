package models

type CommentCreate struct {
	UserID  int    `json:"user_id"`
	PhotoID int    `json:"photo_id"`
	Comment string `json:"comment"`
}
