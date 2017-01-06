package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/buger/jsonparser"

	"strconv"

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

func ListCommentsHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// get query param offset
	var offsetString = r.URL.Query().Get("offset")

	// get query param nr_of_rows
	var rowsString = r.URL.Query().Get("rows")

	// get PhotoID
	var photoIDString = r.URL.Query().Get("photoID")
	photoID, err := strconv.Atoi(photoIDString)
	if err != nil {
		util.SendErrorMessage(w, "photoID is not a number ["+photoIDString+"]")
		return
	}

	dab := db.InitMariaDB()

	comments := make([]*sharedModels.CommentResponse, 0)

	// if not query params
	if len(offsetString) > 0 && len(rowsString) > 0 {
		// list the 10 past on start-lengths
		offset, err := strconv.Atoi(offsetString)
		if err != nil {
			util.SendErrorMessage(w, "offset is not a number ["+offsetString+"]")
			return
		}
		rows, err := strconv.Atoi(rowsString)
		if err != nil {
			util.SendErrorMessage(w, "rows is not a number ["+rowsString+"]")
			return
		}
		comments, err = dab.GetComments(photoID, offset, rows)
	} else {
		// then list last 10
		comments, err = dab.GetComments(photoID, 1, 10)
	}

	if err != nil {
		util.SendError(w, err)
		return
	}

	// return
	util.SendOK(w, comments)
}

func GetCommentCountHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	objects := make([]*sharedModels.CommentCountRequest, 0)
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		comment := &sharedModels.CommentCountRequest{}
		json.Unmarshal(value, comment)
		objects = append(objects, comment)
	}, "requests")

	dab := db.InitMariaDB()
	responses, err := dab.GetCommentCount(objects)

	for i := 0; i < len(responses); i++ {
		log.Printf("Number of coments : %v.\n", responses[i])
	}

	if err != nil {
		util.SendError(w, err)
		return
	}

	type Resp struct {
		Result []*sharedModels.CommentCountResponse `json:"result"`
	}
	util.SendOK(w, &Resp{Result: responses})
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
