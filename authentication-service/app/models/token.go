package models

// Token contains the token properties
type Token struct {
	Token     string `json:"token"`
	ExpiresOn string `json:"expires_on"`
	User      User   `json:"user"`
}
