package middleware

import (
	"errors"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/prometheus/common/log"
	"github.com/urfave/negroni"
)

func AccessControlHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if next != nil {
		next(w, r)
	}
}

func AcceptOPTIONS(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

// RequireTokenAuthenticationHandler is a middleware handler which extracts the token from the header of from the query parameter and checks if the token is valid.
func RequireTokenAuthenticationHandler(secretKey string) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			util.SendBadRequest(w, errors.New("token is mandatory"))
			return
		}

		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			log.Errorf("Error. Token: %v. Message: %v.\n", queryToken, err.Error())
			util.SendBadRequest(w, errors.New("Invalid token"))
			return
		}

		if tok != nil && tok.Valid {
			if next != nil {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(""))
		}
	})
}
