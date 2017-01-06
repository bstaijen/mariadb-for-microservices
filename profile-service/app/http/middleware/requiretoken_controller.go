package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	jwt "github.com/dgrijalva/jwt-go"
)

func RequireTokenAuthenticationController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var queryToken = r.URL.Query().Get("token")

	if len(queryToken) < 1 {
		queryToken = r.Header.Get("token")
	}

	if len(queryToken) < 1 {
		util.SendBadRequest(w, errors.New("token is mandatory"))
		return
	}

	cnf := config.LoadConfig()
	secretKey := cnf.SecretKey
	tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Errorf("Error: %v \n", err.Error())
	}

	if tok != nil && tok.Valid {
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(""))
	}
}
