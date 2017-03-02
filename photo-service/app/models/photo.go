package models

import (
	"time"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// Photo can be used for passing around a photo object in the application.
type Photo struct {
	ID            int                             `json:"id"`
	UserID        int                             `json:"user_id"`
	Username      string                          `json:"username"`
	Title         string                          `json:"title"`
	Filename      string                          `json:"filename"`
	CreatedAt     time.Time                       `json:"createdAt"`
	TotalVotes    int                             `json:"totalVotes"`
	UpvoteCount   int                             `json:"upvote_count"`
	DownvoteCount int                             `json:"downvote_count"`
	YouUpvote     bool                            `json:"upvote"`
	YouDownvote   bool                            `json:"downvote"`
	Comments      []*sharedModels.CommentResponse `json:"comments"`
	CommentCount  int                             `json:"comment_count"`
	ContentType   string                          `json:"-"`
	Image         []byte                          `json:"-"`
}

// CreatePhoto can be used for creating a new photo object
type CreatePhoto struct {
	UserID      int
	Filename    string
	Title       string
	ContentType string
	Image       []byte
}
