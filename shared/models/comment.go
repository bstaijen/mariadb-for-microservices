package models

import "time"

type CommentRequest struct {
	PhotoID int `json:"photo_id"`
}

type CommentResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	PhotoID   int       `json:"photo_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
