package controllers

import (
	"net/http"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

func CreateUserController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := &models.User{}
	err := util.RequestToJSON(r, user)
	if err != nil {
		util.SendBadRequest(w, errors.New("Bad json"))
		return
	}

	if err := user.Validate(); err == nil {
		if err := user.ValidatePassword(); err == nil {

			database := db.InitMariaDB()

			createdID, err := database.CreateUser(user)

			if err != nil {
				util.SendBadRequest(w, err)
				return
			}
			createdUser, _ := database.GetUserByID(createdID)
			util.SendOK(w, createdUser)

		} else {
			util.SendBadRequest(w, err)
		}
	} else {
		util.SendBadRequest(w, err)
	}
}
