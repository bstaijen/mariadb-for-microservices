package util

import (
	"net/http"

	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

func GetUserIDFromRequest(req *http.Request) (int, error) {
	var queryToken = req.URL.Query().Get("token")

	if len(queryToken) < 1 {
		queryToken = req.Header.Get("token")
	}

	if len(queryToken) < 1 {
		return 0, errors.New("No JWT available")
	}

	tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), nil
	})
	if err != nil {
		return 0, err
	}

	claims := tok.Claims.(jwt.MapClaims)
	var ID = claims["sub"].(float64) // gets the ID

	return int(ID), nil
}
