package controllers

import (
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/config"
	"github.com/bstaijen/mariadb-for-microservices/vote-service/database"
	jwt "github.com/dgrijalva/jwt-go"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

func CreateHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var queryToken = r.URL.Query().Get("token")

	if len(queryToken) < 1 {
		queryToken = r.Header.Get("token")
	}

	if len(queryToken) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string("token is mandatory")))
		return
	}

	secretKey := config.LoadConfig().SecretKey
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
		database := db.InitMariaDB()
		err := database.Create(voteCreateObject)
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}

		// 3.send result to frontend
		util.SendOKMessage(w, "You voted")
	} else {
		util.SendErrorMessage(w, "UserID or PhotoID are invalid")
	}

}
