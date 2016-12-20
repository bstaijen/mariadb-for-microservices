package controllers

import (
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/database"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

func CreateHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 1.parse body
	voteCreateObject := &sharedModels.VoteCreateRequest{}

	err := util.RequestToJSON(r, voteCreateObject)
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
		util.SendOKMessage(w, "alright!")
	} else {
		util.SendErrorMessage(w, "UserID or PhotoID are invalid")
	}

}
