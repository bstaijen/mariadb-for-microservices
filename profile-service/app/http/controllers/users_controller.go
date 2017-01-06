package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"strconv"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
)

func UsersController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	database := db.InitMariaDB()
	users, err := database.GetUsers()
	if err != nil {
		util.SendError(w, err)
	} else {
		util.SendOK(w, users)
	}
}

func UserIndexController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	strID := vars["id"]

	id, _ := strconv.Atoi(strID)

	database := db.InitMariaDB()
	user, err := database.GetUserByID(id)

	if err != nil {
		util.SendBadRequest(w, err)
		return
	}
	util.SendOK(w, user)
}

func GetUsernames(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	result := bodyToArrayWithIDs(r)

	database := db.InitMariaDB()
	users, err := database.GetUsernames(result)
	if err != nil {
		util.SendError(w, err)
		return
	}

	type Resp struct {
		Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
	}
	util.SendOK(w, &Resp{Usernames: users})
}

func bodyToArrayWithIDs(req *http.Request) []*sharedModels.GetUsernamesRequest {
	data, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	objects := make([]*sharedModels.GetUsernamesRequest, 0)
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		usernameReq := &sharedModels.GetUsernamesRequest{}
		json.Unmarshal(value, usernameReq)
		objects = append(objects, usernameReq)
	}, "requests")
	return objects
}
