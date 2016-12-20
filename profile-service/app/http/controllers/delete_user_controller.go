package controllers

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

func DeleteUserController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var queryToken = r.URL.Query().Get("token")

	if len(queryToken) < 1 {
		queryToken = r.Header.Get("token")
	}

	if len(queryToken) < 1 {
		util.SendBadRequest(w, errors.New("token is mandatory"))
		return
	}

	user := &models.User{}
	err := util.RequestToJSON(r, user)
	if err != nil {
		util.SendBadRequest(w, errors.New("Bad json"))
		return
	}

	secretKey := config.LoadConfig().SecretKey
	tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	claims := tok.Claims.(jwt.MapClaims)
	var ID = claims["sub"].(float64)

	if int(ID) != user.ID {
		util.SendBadRequest(w, errors.New("you can only delete your own user object"))
		return
	}

	database := db.InitMariaDB()
	database.DeleteUser(user)

	util.SendOK(w, string(""))
}
