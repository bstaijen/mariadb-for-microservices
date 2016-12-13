package models

type VoteCreateRequest struct {
	UserID   int  `json:"user_id"`
	PhotoID  int  `json:"photo_id"`
	Upvote   bool `json:upvote`
	Downvote bool `json:downvote`
}

type HasVotedRequest struct {
	UserID  int `json:"user_id"`
	PhotoID int `json:"photo_id"`
}
type HasVotedResponse struct {
	UserID   int  `json:"user_id"`
	PhotoID  int  `json:"photo_id"`
	Downvote bool `json:"downvote"`
	Upvote   bool `json:"upvote"`
}

type VoteCountRequest struct {
	PhotoID int `json:"photo_id"`
}
type VoteCountResponse struct {
	PhotoID int `json:"photo_id"`
	Count   int `json:"count"`
}

type TopRatedPhotoResponse struct {
	PhotoID int `json:"photo_id"`
}
