package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/urfave/negroni"

	"strconv"

	"github.com/bstaijen/mariadb-for-microservices/shared/helper"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

func CreateHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		commentObject := &models.CommentCreate{}
		err := util.RequestToJSON(r, commentObject)
		if err != nil {
			util.SendErrorMessage(w, "bad json")
			return
		}
		if commentObject.UserID > 0 && commentObject.PhotoID > 0 && len(commentObject.Comment) > 0 {
			comment, err := db.Create(connection, commentObject)
			if err != nil {
				util.SendBadRequest(w, err)
				return
			}
			identifiers := make([]*sharedModels.GetUsernamesRequest, 0)
			identifiers = append(identifiers, &sharedModels.GetUsernamesRequest{
				ID: comment.UserID,
			})
			username := getUsernames(cnf, identifiers)
			if len(username) > 0 {
				comment.Username = username[0].Username
			}
			util.SendOK(w, comment)

			return
		}
		util.SendErrorMessage(w, "UserID, PhotoID and Comment are mandatory")
	})
}

func ListCommentsHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		offset, rows := helper.PaginationFromRequest(r)

		// get PhotoID
		var photoIDString = r.URL.Query().Get("photoID")
		photoID, err := strconv.Atoi(photoIDString)
		if err != nil {
			util.SendErrorMessage(w, "photoID is not a number ["+photoIDString+"]")
			return
		}

		comments := make([]*sharedModels.CommentResponse, 0)

		comments, err = db.GetComments(connection, photoID, offset, rows)
		if err != nil {
			util.SendError(w, err)
			return
		}

		// include usernames
		identifiers := make([]*sharedModels.GetUsernamesRequest, 0)
		for index := 0; index < len(comments); index++ {
			comment := comments[index]
			identifiers = append(identifiers, &sharedModels.GetUsernamesRequest{
				ID: comment.UserID,
			})
		}
		usernames := getUsernames(cnf, identifiers)
		for index := 0; index < len(comments); index++ {
			comment := comments[index]
			for undex := 0; undex < len(usernames); undex++ {
				username := usernames[undex]

				if comment.UserID == username.ID {
					comment.Username = username.Username
				}
			}
		}

		// return
		util.SendOK(w, comments)
	})
}

func GetCommentCountHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.Body == nil {
			util.SendErrorMessage(w, "bad json")
			return
		}

		type Collection struct {
			Objects []*sharedModels.CommentCountRequest `json:"requests"`
		}
		col := &Collection{}
		col.Objects = make([]*sharedModels.CommentCountRequest, 0)
		err := util.RequestToJSON(r, &col)
		if err != nil {
			log.Fatal(err)
		}

		responses, err := db.GetCommentCount(connection, col.Objects)

		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Result []*sharedModels.CommentCountResponse `json:"result"`
		}
		util.SendOK(w, &Resp{Result: responses})
	})
}

func GetLastTenHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		type Collection struct {
			Objects []*sharedModels.CommentRequest `json:"requests"`
		}
		col := &Collection{}
		col.Objects = make([]*sharedModels.CommentRequest, 0)
		err := util.RequestToJSON(r, &col)
		if err != nil {
			util.SendError(w, err)
			return
		}

		responses, err := db.GetLastTenComments(connection, col.Objects)
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
		usernames := getUsernames(cnf, identifiers)
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
	})
}

func getUsernames(cnf config.Config, input []*sharedModels.GetUsernamesRequest) []*sharedModels.GetUsernamesResponse {
	type Req struct {
		Requests []*sharedModels.GetUsernamesRequest `json:"requests"`
	}
	body, _ := json.Marshal(Req{Requests: input})
	// Make url
	url := cnf.ProfileServiceBaseurl + "ipc/usernames"

	// Return object
	usernames := make([]*sharedModels.GetUsernamesResponse, 0)

	// GET data and append to return object
	util.Request("GET", url, body, func(res *http.Response) {
		// Error handling
		if res.StatusCode < 200 || res.StatusCode > 299 {
			logrus.Errorf("ERROR with code %v.", res.Status)
			data, _ := ioutil.ReadAll(res.Body)
			logrus.Error(string(data))
			return
		}

		type Collection struct {
			Objects []*sharedModels.GetUsernamesResponse `json:"usernames"`
		}
		col := &Collection{}
		col.Objects = make([]*sharedModels.GetUsernamesResponse, 0)

		err := util.ResponseJSONToObject(res, &col)
		if err != nil {
			log.Fatal(err)
			return
		}
		usernames = col.Objects
	})
	return usernames
}
