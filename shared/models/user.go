package models

type GetUsernamesRequest struct {
	ID int `json:"id"`
}

type GetUsernamesResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
