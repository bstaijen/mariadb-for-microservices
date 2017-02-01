package models

import "time"

// CommentRequest is a struct containing the fields needed for the get comments IPC.
type CommentRequest struct {
	PhotoID int `json:"photo_id"`
}

// CommentResponse is a struct containing the fields the comments IPC returns
type CommentResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	PhotoID   int       `json:"photo_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

// CommentCountRequest is a struct containing the fields needed for the CommentCount IPC
type CommentCountRequest struct {
	PhotoID int `json:"photo_id"`
}

// CommentCountResponse is a struct containing the fields the CommentCount IPC returns
type CommentCountResponse struct {
	PhotoID int `json:"photo_id"`
	Count   int `json:"count"`
}
