package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"strconv"

	"github.com/bstaijen/mariadb-for-microservices/shared/helper"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// CreateHandler creates a comment and stores it in the database
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

// ListCommentsHandler return a list of comments
func ListCommentsHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logrus.Info("List comments")
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

// ListCommentsFromUser : Return all comments from an user and add the photo(extra information) too.
func ListCommentsFromUser(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logrus.Info("List comments from user")

		offset, rows := helper.PaginationFromRequest(r)

		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(string("token is mandatory")))
			return
		}

		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(cnf.SecretKey), nil
		})

		if err != nil {
			util.SendErrorMessage(w, "You are not authorized")
			return
		}

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64) // gets the ID

		comments, err := db.GetCommentsByUserID(connection, int(ID), offset, rows)

		// collect photo IDs.
		ids := make([]*sharedModels.TopRatedPhotoResponse, 0)
		f := make([]*sharedModels.HasVotedRequest, 0)
		g := make([]*sharedModels.VoteCountRequest, 0)
		for _, v := range comments {
			ids = append(ids, &sharedModels.TopRatedPhotoResponse{
				PhotoID: v.PhotoID,
			})

			f = append(f, &sharedModels.HasVotedRequest{
				PhotoID: v.PhotoID,
				UserID:  v.UserID,
			})

			g = append(g, &sharedModels.VoteCountRequest{
				PhotoID: v.PhotoID,
			})
		}

		// get photos
		photos := getPhotos(cnf, ids)

		// get votes
		photos = appendUserVoted(cnf, f, photos)
		photos = appendVotesCount(cnf, g, photos)

		// merge with comments
		type Res struct {
			Comment *sharedModels.CommentResponse `json:"comment"`
			Photo   *sharedModels.PhotoResponse   `json:"photo"`
		}

		h := make([]*Res, 0)
		for _, v := range comments {

			for _, photo := range photos {
				if v.PhotoID == photo.ID {
					h = append(h, &Res{
						Comment: v,
						Photo:   photo,
					})
				}
			}
		}

		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Result []*Res `json:"result"`
		}
		util.SendOK(w, &Resp{Result: h})
	})
}

// GetCommentCountHandler returns a list of counts beloning to comments.
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

// GetLastTenHandler returns a list of last 10 comments
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

// DeleteCommentHandler : is the handler to remove a comment in the database
func DeleteCommentHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
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

		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(cnf.SecretKey), nil
		})

		if err != nil {
			util.SendErrorMessage(w, "You are not authorized")
			return
		}

		claims := tok.Claims.(jwt.MapClaims)
		var userID = claims["sub"].(float64) // gets the ID

		// Get commentID
		vars := mux.Vars(r)
		strID := vars["id"]
		commentID, err := strconv.Atoi(strID)

		if err != nil {
			util.SendErrorMessage(w, "id needs to be numeric")
			return
		}

		if commentID < 1 {
			util.SendErrorMessage(w, "id needs to be greater than 0")
			return
		}

		comment, err := db.GetCommentByID(connection, commentID)
		if err != nil {
			util.SendError(w, err)
			return
		}

		if comment.UserID != int(userID) {
			util.SendErrorMessage(w, "you can only remove your own comment")
			return
		}

		_, err = db.DeleteCommentByID(connection, commentID)
		if err != nil {
			util.SendError(w, err)
			return
		}
		util.SendOKMessage(w, "Comment removed")
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

func getPhotos(cnf config.Config, input []*sharedModels.TopRatedPhotoResponse) []*sharedModels.PhotoResponse {
	type Req struct {
		Requests []*sharedModels.TopRatedPhotoResponse `json:"requests"`
	}
	body, _ := json.Marshal(&Req{Requests: input})

	// Make url
	url := cnf.PhotoServiceBaseurl + "ipc/getPhotos"

	photos := make([]*sharedModels.PhotoResponse, 0)
	if strings.HasPrefix(url, "http") {
		err := util.Request("GET", url, body, func(res *http.Response) {
			// Error handling
			if res.StatusCode < 200 || res.StatusCode > 299 {
				printResponseError(res)
				return
			}

			// Happy path
			type Collection struct {
				Objects []*sharedModels.PhotoResponse `json:"results"`
			}
			col := &Collection{}
			col.Objects = make([]*sharedModels.PhotoResponse, 0)
			err := util.ResponseJSONToObject(res, &col)
			if err != nil {
				log.Fatal(err)
			}
			photos = col.Objects
		})
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Errorf("Wrong URL. Expected something which starts with http, instead got %v.", url)
	}
	return photos
}

// appendVotesCount triggers `GET votes request` and appends result to []Photo
func appendVotesCount(cnf config.Config, photoCountIdentifiers []*sharedModels.VoteCountRequest, photos []*sharedModels.PhotoResponse) []*sharedModels.PhotoResponse {
	// Get the vote counts and add them
	results := getVotes(cnf, photoCountIdentifiers)
	for index := 0; index < len(photos); index++ {
		photoObject := photos[index]

		for resultsIndex := 0; resultsIndex < len(results); resultsIndex++ {
			resultsObject := results[resultsIndex]
			if photoObject.ID == resultsObject.PhotoID {
				photoObject.TotalVotes = resultsObject.UpVoteCount + resultsObject.DownVoteCount
				photoObject.UpvoteCount = resultsObject.UpVoteCount
				photoObject.DownvoteCount = resultsObject.DownVoteCount
			}
		}
	}
	return photos
}

// appendUserVoted triggers `GET voted request` and appends results to []Photo. This function will
// lookup whether the user has voted on the Photo
func appendUserVoted(cnf config.Config, photoVotedIdentifiers []*sharedModels.HasVotedRequest, photos []*sharedModels.PhotoResponse) []*sharedModels.PhotoResponse {
	youVoted := voted(cnf, photoVotedIdentifiers)

	for index := 0; index < len(photos); index++ {
		photoObject := photos[index]
		for votesIndex := 0; votesIndex < len(youVoted); votesIndex++ {
			obj := youVoted[votesIndex]

			if photoObject.ID == obj.PhotoID {
				photoObject.YouUpvote = obj.Upvote
				photoObject.YouDownvote = obj.Downvote
			}
		}
	}
	return photos
}

// Get votes from the VotesSerivce
func getVotes(cnf config.Config, input []*sharedModels.VoteCountRequest) []*sharedModels.VoteCountResponse {
	type Req struct {
		Requests []*sharedModels.VoteCountRequest `json:"requests"`
	}
	body, _ := json.Marshal(&Req{Requests: input})

	// Make url
	url := cnf.VoteServiceBaseurl + "ipc/count"

	//Return object
	votes := make([]*sharedModels.VoteCountResponse, 0)
	if strings.HasPrefix(url, "http") {
		err := util.Request("GET", url, body, func(res *http.Response) {
			// Error handling
			if res.StatusCode < 200 || res.StatusCode > 299 {
				printResponseError(res)
				return
			}

			// Happy path
			type Collection struct {
				Objects []*sharedModels.VoteCountResponse `json:"results"`
			}
			col := &Collection{}
			col.Objects = make([]*sharedModels.VoteCountResponse, 0)
			err := util.ResponseJSONToObject(res, &col)
			if err != nil {
				logrus.Warn(err)
			}
			votes = col.Objects
		})
		if err != nil {
			logrus.Warn(err)
		}
	} else {
		logrus.Errorf("Wrong URL. Expected something which starts with http, instead got %v.", url)
	}
	return votes
}

// Determine if the user has voted on a photo. VotesService
func voted(cnf config.Config, input []*sharedModels.HasVotedRequest) []*sharedModels.HasVotedResponse {
	type Req struct {
		Requests []*sharedModels.HasVotedRequest `json:"requests"`
	}
	body, _ := json.Marshal(&Req{Requests: input})

	// Make url
	url := cnf.VoteServiceBaseurl + "ipc/voted"

	//Return object
	hasVoted := make([]*sharedModels.HasVotedResponse, 0)
	if strings.HasPrefix(url, "http") {
		err := util.Request("GET", url, body, func(res *http.Response) {
			// Error handling
			if res.StatusCode < 200 || res.StatusCode > 299 {
				printResponseError(res)
				return
			}

			// Happy path
			type Collection struct {
				Objects []*sharedModels.HasVotedResponse `json:"results"`
			}
			col := &Collection{}
			col.Objects = make([]*sharedModels.HasVotedResponse, 0)
			err := util.ResponseJSONToObject(res, &col)
			if err != nil {
				logrus.Warn(err)
			}
			hasVoted = col.Objects
		})
		if err != nil {
			logrus.Warn(err)
		}
	} else {
		logrus.Errorf("Wrong URL. Expected something which starts with http, instead got %v.", url)
	}
	return hasVoted
}

func printResponseError(res *http.Response) {
	logrus.Errorf("Response error with statuscode %v.", res.Status)
	data, _ := ioutil.ReadAll(res.Body)
	logrus.Error(string(data))
}
