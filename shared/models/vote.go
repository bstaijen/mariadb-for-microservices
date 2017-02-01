package models

// VoteCreateRequest contains all fields needed for a new vote.
type VoteCreateRequest struct {
	UserID   int  `json:"user_id"`
	PhotoID  int  `json:"photo_id"`
	Upvote   bool `json:"upvote"`
	Downvote bool `json:"downvote"`
}

// HasVotedRequest contains all fields needed for the HasVoted IPC.
type HasVotedRequest struct {
	UserID  int `json:"user_id"`
	PhotoID int `json:"photo_id"`
}

// HasVotedResponse contains the fields the HasVoted IPC returns
type HasVotedResponse struct {
	UserID   int  `json:"user_id"`
	PhotoID  int  `json:"photo_id"`
	Downvote bool `json:"downvote"`
	Upvote   bool `json:"upvote"`
}

// VoteCountRequest contains all fields needed for the VoteCount IPC
type VoteCountRequest struct {
	PhotoID int `json:"photo_id"`
}

// VoteCountResponse contains all the fields the VoteCount IPC returns
type VoteCountResponse struct {
	PhotoID       int `json:"photo_id"`
	UpVoteCount   int `json:"total_up_count"`
	DownVoteCount int `json:"total_down_count"`
}

// TopRatedPhotoResponse contains the field the TopRated IPC returns
type TopRatedPhotoResponse struct {
	PhotoID int `json:"photo_id"`
}
