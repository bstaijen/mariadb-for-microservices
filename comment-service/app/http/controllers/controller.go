package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/buger/jsonparser"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

func CreateHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	commentObject := &models.CommentCreate{}
	err := util.RequestToJSON(r, commentObject)
	if err != nil {
		util.SendErrorMessage(w, "bad json")
		return
	}
	if commentObject.UserID > 0 && commentObject.PhotoID > 0 && len(commentObject.Comment) > 0 {
		database := db.InitMariaDB()
		comment, err := database.Create(commentObject)
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}
		identifiers := make([]*sharedModels.GetUsernamesRequest, 0)
		identifiers = append(identifiers, &sharedModels.GetUsernamesRequest{
			ID: comment.UserID,
		})
		username := getUsernames(identifiers)
		if len(username) > 0 {
			comment.Username = username[0].Username
		}
		util.SendOK(w, comment)

		return
	}
	util.SendErrorMessage(w, "UserID, PhotoID and Comment are mandatory")
}

func GetLastTenHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	objects := make([]*sharedModels.CommentRequest, 0)
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		comment := &sharedModels.CommentRequest{}
		json.Unmarshal(value, comment)
		objects = append(objects, comment)
	}, "requests")

	dab := db.InitMariaDB()
	responses, err := dab.GetLastTenComments(objects)
	if err != nil {
		util.SendError(w, err)
		return
	}

	// include usernames
	identifiers := make([]*sharedModels.GetUsernamesRequest, 0)
	for index := 0; index < len(responses); index++ {
		comment := responses[index]
		identifiers = append(identifiers, &sharedModels.GetUsernamesRequest{
			ID: comment.UserID,
		})
	}
	usernames := getUsernames(identifiers)
	for index := 0; index < len(responses); index++ {
		comment := responses[index]
		for undex := 0; undex < len(usernames); undex++ {
			username := usernames[undex]

			if comment.UserID == username.ID {
				comment.Username = username.Username
			}
		}
	}

	type Resp struct {
		Comments []*sharedModels.CommentResponse `json:"comments"`
	}
	util.SendOK(w, &Resp{Comments: responses})
}

func getUsernames(input []*sharedModels.GetUsernamesRequest) []*sharedModels.GetUsernamesResponse {
	type Req struct {
		Requests []*sharedModels.GetUsernamesRequest `json:"requests"`
	}
	body, _ := json.Marshal(Req{Requests: input})
	// Make url
	url := config.LoadConfig().ProfileServiceBaseurl + "ipc/usernames"

	// Return object
	usernames := make([]*sharedModels.GetUsernamesResponse, 0)

	// GET data and append to return object
	util.Request("GET", url, body, func(res *http.Response) {
		data, err := ioutil.ReadAll(res.Body)
		if err == nil {
			jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				username := &sharedModels.GetUsernamesResponse{}
				json.Unmarshal(value, username)
				usernames = append(usernames, username)
			}, "usernames")
		}
	})
	return usernames
}
