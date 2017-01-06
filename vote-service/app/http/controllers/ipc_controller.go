package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/database"

	"io/ioutil"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/buger/jsonparser"
)

func GetTopRatedHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	database := db.InitMariaDB()
	topRated, err := database.GetTopRatedTimeline()
	if err != nil {
		util.SendError(w, err)
		return
	}

	type Resp struct {
		Results []*sharedModels.TopRatedPhotoResponse `json:"results"`
	}
	util.SendOK(w, &Resp{Results: topRated})
}

func GetHotHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	database := db.InitMariaDB()
	hot, err := database.GetHotTimeline()
	if err != nil {
		util.SendError(w, err)
		return
	}

	type Resp struct {
		Results []*sharedModels.TopRatedPhotoResponse `json:"results"`
	}
	util.SendOK(w, &Resp{Results: hot})
}

func HasVotedHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	objects := make([]*sharedModels.HasVotedRequest, 0)
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		vote := &sharedModels.HasVotedRequest{}
		json.Unmarshal(value, vote)
		objects = append(objects, vote)
	}, "requests")

	database := db.InitMariaDB()
	counts, err := database.HasVoted(objects)
	if err != nil {
		util.SendError(w, err)
		return
	}

	type Resp struct {
		Results []*sharedModels.HasVotedResponse `json:"results"`
	}
	util.SendOK(w, &Resp{Results: counts})
}

func GetVoteCountHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	objects := make([]*sharedModels.VoteCountRequest, 0)
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		vote := &sharedModels.VoteCountRequest{}
		json.Unmarshal(value, vote)
		objects = append(objects, vote)
	}, "requests")

	database := db.InitMariaDB()
	counts, err := database.VoteCount(objects)
	if err != nil {
		util.SendError(w, err)
		return
	}

	type Resp struct {
		Results []*sharedModels.VoteCountResponse `json:"results"`
	}
	util.SendOK(w, &Resp{Results: counts})
}
