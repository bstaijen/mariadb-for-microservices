package controllers

import (
	"database/sql"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

func UpdateUserHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(string("token is mandatory")))
			return
		}

		user := &models.User{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("bad json"))
			return
		}

		secretKey := cnf.SecretKey
		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64) // gets the ID

		if int(ID) != user.ID {
			util.SendBadRequest(w, errors.New("you can only change your own user object"))
			return
		}

		if err := user.Validate(); err == nil {
			if err := user.ValidatePassword(); err == nil {

				db.UpdateUser(connection, user)

				util.SendOK(w, user)

			} else {
				util.SendBadRequest(w, err)
			}
		} else {
			util.SendBadRequest(w, err)
		}
	})
}
