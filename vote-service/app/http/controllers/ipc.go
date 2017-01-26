package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/database"
	"github.com/urfave/negroni"

	"io/ioutil"

	"github.com/bstaijen/mariadb-for-microservices/shared/helper"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/buger/jsonparser"
)

func GetTopRatedHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		// get offset and rows
		offset, rows := helper.PaginationFromRequest(r)

		topRated, err := db.GetTopRatedTimeline(connection, offset, rows)
		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Results []*sharedModels.TopRatedPhotoResponse `json:"results"`
		}
		util.SendOK(w, &Resp{Results: topRated})
	})
}

func GetHotHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		// get offset and rows
		offset, rows := helper.PaginationFromRequest(r)

		hot, err := db.GetHotTimeline(connection, offset, rows)
		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Results []*sharedModels.TopRatedPhotoResponse `json:"results"`
		}
		util.SendOK(w, &Resp{Results: hot})
	})
}

func HasVotedHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		data, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		objects := make([]*sharedModels.HasVotedRequest, 0)
		jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			vote := &sharedModels.HasVotedRequest{}
			json.Unmarshal(value, vote)
			objects = append(objects, vote)
		}, "requests")

		counts, err := db.HasVoted(connection, objects)
		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Results []*sharedModels.HasVotedResponse `json:"results"`
		}
		util.SendOK(w, &Resp{Results: counts})
	})
}

func GetVoteCountHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		data, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		objects := make([]*sharedModels.VoteCountRequest, 0)
		jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			vote := &sharedModels.VoteCountRequest{}
			json.Unmarshal(value, vote)
			objects = append(objects, vote)
		}, "requests")

		counts, err := db.VoteCount(connection, objects)
		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Results []*sharedModels.VoteCountResponse `json:"results"`
		}
		util.SendOK(w, &Resp{Results: counts})
	})
}
