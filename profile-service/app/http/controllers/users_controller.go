package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func UserByIndexHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		vars := mux.Vars(r)
		strID := vars["id"]

		id, err := strconv.Atoi(strID)
		if err != nil {
			logrus.Error(err)
		}

		user, err := db.GetUserByID(connection, id)

		if err != nil {
			util.SendBadRequest(w, err)
			return
		}
		util.SendOK(w, user)
	})
}

func GetUsernamesHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result, err := bodyToArrayWithIDs(r)

		if err != nil {
			util.SendError(w, err)
			return
		}

		users, err := db.GetUsernames(connection, result)
		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
		}
		util.SendOK(w, &Resp{Usernames: users})
	})
}

func bodyToArrayWithIDs(req *http.Request) ([]*sharedModels.GetUsernamesRequest, error) {
	data, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	objects := make([]*sharedModels.GetUsernamesRequest, 0)
	err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		usernameReq := &sharedModels.GetUsernamesRequest{}
		json.Unmarshal(value, usernameReq)
		objects = append(objects, usernameReq)
	}, "requests")
	if err != nil {
		return nil, err
	}
	return objects, nil
}
