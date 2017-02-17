package models

import (
	"time"
)

type PhotoRequest struct {
	PhotoID int `json:"photo_id"`
}

type PhotoResponse struct {
	ID            int                `json:"id"`
	UserID        int                `json:"user_id"`
	Username      string             `json:"username"`
	Title         string             `json:"title"`
	Filename      string             `json:"filename"`
	CreatedAt     time.Time          `json:"createdAt"`
	TotalVotes    int                `json:"totalVotes"`
	UpvoteCount   int                `json:"upvote_count"`
	DownvoteCount int                `json:"downvote_count"`
	YouUpvote     bool               `json:"upvote"`
	YouDownvote   bool               `json:"downvote"`
	Comments      []*CommentResponse `json:"comments"`
	CommentCount  int                `json:"comment_count"`
}
