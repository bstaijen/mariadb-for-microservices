package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"
)

// AccessControlHandler sets the Access-Control-Allow-Origin header on a request header. The header specifies the URI that may access the resources.
func AccessControlHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if next != nil {
		next(w, r)
	}
}

// AcceptOPTIONSHandler sets the Access-Control-Allow-Origin header on a request header. The header specifies the URI that may access the resources. AcceptOPTIONSHandler also sets the Access-Control-Allow-Headers which lets the server whitelist headers that browsers are allowed to access.
func AcceptOPTIONSHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

// RequireTokenAuthenticationHandler is a middleware handler which extracts the token from the header of from the query parameter and checks if the token is valid.
func RequireTokenAuthenticationHandler(cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			util.SendBadRequest(w, errors.New("token is mandatory"))
			return
		}

		secretKey := cnf.SecretKey
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
