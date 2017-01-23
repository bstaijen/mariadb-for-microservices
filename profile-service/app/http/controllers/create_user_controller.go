package controllers

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/urfave/negroni"

	log "github.com/Sirupsen/logrus"
)

func CreateUserHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := &models.User{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("Bad json"))
			return
		}

		if err := user.Validate(); err == nil {
			if err := user.ValidatePassword(); err == nil {

				hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
				user.Hash = string(hash)

				createdID, err := db.CreateUser(connection, user)

				if err != nil {
					util.SendBadRequest(w, err)
					return
				}
				log.Debugln("GetUserByID")
				createdUser, _ := db.GetUserByID(connection, createdID)
				util.SendOK(w, createdUser)

			} else {
				util.SendBadRequest(w, err)
			}
		} else {
			util.SendBadRequest(w, err)
		}
	})
}
