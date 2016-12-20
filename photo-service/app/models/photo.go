package models

import (
	"time"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// Photo can be used for passing around a photo object in the application.
type Photo struct {
	ID          int                             `json:"id"`
	UserID      int                             `json:"user_id"`
	Username    string                          `json:"username"`
	Title       string                          `json:"title"`
	Filename    string                          `json:"filename"`
	CreatedAt   time.Time                       `json:"createdAt"`
	TotalVotes  int                             `json:"totalVotes"`
	YouUpvote   bool                            `json:"upvote"`
	YouDownvote bool                            `json:"downvote"`
	Comments    []*sharedModels.CommentResponse `json:"comments"`

	ContentType string
	Image       []byte
}

// CreatePhoto can be used for creating a new photo object
type CreatePhoto struct {
	UserID      int
	Filename    string
	Title       string
	ContentType string
	Image       []byte
}
