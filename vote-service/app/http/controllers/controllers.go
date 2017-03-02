package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/vote-service/config"
	"github.com/bstaijen/mariadb-for-microservices/vote-service/database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"

	"github.com/bstaijen/mariadb-for-microservices/shared/helper"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

// CreateHandler handler for creating a vote and storing it in the database.
func CreateHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
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

		secretKey := cnf.SecretKey
		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			util.SendErrorMessage(w, "You are not authorized")
			return
		}

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64) // gets the ID

		// 1.parse body
		voteCreateObject := &sharedModels.VoteCreateRequest{}
		voteCreateObject.UserID = int(ID)

		err = util.RequestToJSON(r, voteCreateObject)
		if err != nil {
			util.SendErrorMessage(w, "bad json")
			return
		}

		if voteCreateObject.Upvote == voteCreateObject.Downvote {
			util.SendErrorMessage(w, "can not vote for none or both")
			return
		}

		if voteCreateObject.UserID > 0 && voteCreateObject.PhotoID > 0 && (voteCreateObject.Upvote || voteCreateObject.Downvote) {
			// 2.save in database
			err := db.Create(connection, voteCreateObject)
			if err != nil {
				util.SendBadRequest(w, err)
				return
			}

			// 3.send result to frontend
			util.SendOKMessage(w, "You voted")
		} else {
			util.SendErrorMessage(w, "UserID or PhotoID are invalid")
		}

	})
}

func GetVotesFromAUser(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
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

		secretKey := cnf.SecretKey
		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			util.SendErrorMessage(w, "You are not authorized")
			return
		}

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64) // gets the ID

		// get offset and rows
		offset, rows := helper.PaginationFromRequest(r)

		photoIDs, err := db.GetVotesFromUser(connection, int(ID), offset, rows)
		if err != nil {
			util.SendError(w, err)
			return
		}

		photos := getPhotos(cnf, photoIDs)

		t := make([]*sharedModels.HasVotedRequest, 0)
		g := make([]*sharedModels.VoteCountRequest, 0)
		for _, v := range photos {
			t = append(t, &sharedModels.HasVotedRequest{
				PhotoID: v.ID,
				UserID:  v.UserID,
			})
			g = append(g, &sharedModels.VoteCountRequest{
				PhotoID: v.ID,
			})
		}

		// Get youVoted
		photos = appendUserVoted(connection, cnf, t, photos)
		photos = appendVotesCount(connection, cnf, g, photos)

		type Resp struct {
			Result []*sharedModels.PhotoResponse `json:"result"`
		}
		util.SendOK(w, &Resp{Result: photos})
	})
}

func HealthHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		log.Println("Health Handler called")

		util.SendOKMessage(w, "I am healthy")

	})
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
func appendVotesCount(connection *sql.DB, cnf config.Config, photoCountIdentifiers []*sharedModels.VoteCountRequest, photos []*sharedModels.PhotoResponse) []*sharedModels.PhotoResponse {
	// Get the vote counts and add them
	results, err := db.VoteCount(connection, photoCountIdentifiers)
	if err != nil {
		logrus.Fatal(err)
	}
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
func appendUserVoted(connection *sql.DB, cnf config.Config, photoVotedIdentifiers []*sharedModels.HasVotedRequest, photos []*sharedModels.PhotoResponse) []*sharedModels.PhotoResponse {
	youVoted, err := db.HasVoted(connection, photoVotedIdentifiers)
	if err != nil {
		logrus.Fatal(err)
	}
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

func printResponseError(res *http.Response) {
	logrus.Errorf("Response error with statuscode %v.", res.Status)
	data, _ := ioutil.ReadAll(res.Body)
	logrus.Error(string(data))
}
