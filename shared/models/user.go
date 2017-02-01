package models

// GetUsernamesRequest is a struct and contains the fields the GetUsername IPC needs
type GetUsernamesRequest struct {
	ID int `json:"id"`
}

// GetUsernamesResponse is a struct and contains the fields the GetUsername IPC returns
type GetUsernamesResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
