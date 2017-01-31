package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/database"
	"github.com/urfave/negroni"

	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
	"github.com/bstaijen/mariadb-for-microservices/shared/helper"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/buger/jsonparser"
	jwt "github.com/dgrijalva/jwt-go"
)

// IncomingHandler is the handler for serving the default photos timeline
func IncomingHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// Get user ID. It is allowed to be 0.
		userID, _ := getUserIDFromRequest(cnf, r)

		// get offset and rows
		offset, rows := helper.PaginationFromRequest(r)

		photos, err := db.ListIncoming(connection, offset, rows)

		logrus.Infof("Number of photos retrieved from database : %v.", len(photos))

		if err != nil {
			util.SendError(w, err)
			return
		}
		photos = findResources(cnf, photos, userID, true, true, true)

		for _, v := range photos {
			for _, s := range v.Comments {
				logrus.Infof("test %v", s.Comment)
			}
		}

		util.SendOK(w, photos)
	})
}

// TopRatedHandler is the handler for serving the Top Rated photos timeline
func TopRatedHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// Get user ID. It is allowed to be 0.
		userID, _ := getUserIDFromRequest(cnf, r)

		// get offset and rows and pass into URL toprated
		offset, rows := helper.PaginationFromRequest(r)

		// Make url
		urlpart := fmt.Sprintf("ipc/toprated?offset=%v&rows=%v", offset, rows)
		url := cnf.VoteServiceBaseurl + urlpart

		// Save object
		photos := make([]*models.Photo, 0)

		// GET
		util.Request("GET", url, []byte(string("")), func(res *http.Response) {
			data, err := ioutil.ReadAll(res.Body)
			if err == nil {
				jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					obj := &sharedModels.TopRatedPhotoResponse{}
					json.Unmarshal(value, obj)
					photo, err := db.GetPhotoById(connection, obj.PhotoID)
					if photo != nil {
						photos = append(photos, photo)
					}
				}, "results")
			}
		})

		photos = findResources(cnf, photos, userID, true, true, true)
		util.SendOK(w, photos)
	})
}

// HotHandler is the handler for serving the Hot photos timeline
func HotHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		userID, _ := getUserIDFromRequest(cnf, r)

		// get offset and rows and pass into URL toprated
		offset, rows := helper.PaginationFromRequest(r)

		// Make url
		urlpart := fmt.Sprintf("ipc/hot?offset=%v&rows=%v", offset, rows)
		url := cnf.VoteServiceBaseurl + urlpart

		// Save object
		photos := make([]*models.Photo, 0)

		// GET
		util.Request("GET", url, []byte(string("")), func(res *http.Response) {
			// Error handling
			if res.StatusCode < 200 || res.StatusCode > 299 {
				logrus.Errorf("ERROR with code %v.", res.Status)
				data, _ := ioutil.ReadAll(res.Body)
				logrus.Error(string(data))

				util.SendErrorMessage(w, "Could not retrieve photos.")
				return
			}

			// Happy path
			data, err := ioutil.ReadAll(res.Body)
			if err == nil {
				jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					obj := &sharedModels.TopRatedPhotoResponse{}
					json.Unmarshal(value, obj)
					photo, err := db.GetPhotoById(connection, obj.PhotoID)
					photos = append(photos, photo)
				}, "results")
			}
		})

		photos = findResources(cnf, photos, userID, true, true, true)
		util.SendOK(w, photos)
	})
}

// FindResources searches for related resources to a collection of photos and adds them to the photo object.
// By specifying parameters the caller of this func can determine which resources will be added and which
// will be skippd. If userID is 0 or less then this func cannot determine if the user has voted on the photos.
func findResources(cnf config.Config, photos []*models.Photo, userID int, comments bool, usernames bool, votes bool) []*models.Photo {

	// PhotoID holder
	photoCountIdentifiers := make([]*sharedModels.VoteCountRequest, 0)
	photoCommentCountIdentifiers := make([]*sharedModels.CommentCountRequest, 0)
	photoVotedIdentifiers := make([]*sharedModels.HasVotedRequest, 0)
	photoCommentsIdentifiers := make([]*sharedModels.CommentRequest, 0)
	photoUsernamesIndentifiers := make([]*sharedModels.GetUsernamesRequest, 0)

	// Append the username to the photos
	for index := 0; index < len(photos); index++ {
		photoObject := photos[index]

		// Collect photoIDs
		photoCountIdentifiers = append(photoCountIdentifiers, &sharedModels.VoteCountRequest{
			PhotoID: photoObject.ID,
		})
		photoCommentCountIdentifiers = append(photoCommentCountIdentifiers, &sharedModels.CommentCountRequest{
			PhotoID: photoObject.ID,
		})
		photoVotedIdentifiers = append(photoVotedIdentifiers, &sharedModels.HasVotedRequest{
			PhotoID: photoObject.ID,
			UserID:  userID,
		})
		photoCommentsIdentifiers = append(photoCommentsIdentifiers, &sharedModels.CommentRequest{
			PhotoID: photoObject.ID,
		})
		photoUsernamesIndentifiers = append(photoUsernamesIndentifiers, &sharedModels.GetUsernamesRequest{
			ID: photoObject.UserID,
		})

	} // end: for photos

	// Searches and adds comments to the photos
	if comments {
		photos = appendComments(cnf, photoCommentsIdentifiers, photos)
		photos = appendCommentCount(cnf, photoCommentCountIdentifiers, photos)
	}

	// Searches and adds usernames to the photos
	if usernames {
		photos = appendUsernames(cnf, photoUsernamesIndentifiers, photos)
	}

	// Searches and adds the votes count on each photo and
	// whether or not the user has voted on this particular picture.
	if votes {

		photos = appendVotesCount(cnf, photoCountIdentifiers, photos)

		// Get up/downvote from requesting user and add them
		if userID > 0 {
			photos = appendUserVoted(cnf, photoVotedIdentifiers, photos)
		} else {
			fmt.Printf("Warning: UserID is to small for voting. User ID : %v\n", userID)
		}
	}
	return photos
}

// appendComments triggers `GET comments request` (last 10 comments for each photo) and appends results to []Photo
func appendComments(cnf config.Config, photoCommentsIdentifiers []*sharedModels.CommentRequest, photos []*models.Photo) []*models.Photo {
	for _, v := range photoCommentsIdentifiers {
		logrus.Infof("photoCommentsIdentifiers %v", v)
	}

	comments := getComments(cnf, photoCommentsIdentifiers)
	for ind := 0; ind < len(photos); ind++ {

		// Get reference
		phot := photos[ind]
		// Create an array
		phot.Comments = make([]*sharedModels.CommentResponse, 0)

		for i := 0; i < len(comments); i++ {
			commentObject := comments[i]
			if phot.ID == commentObject.PhotoID {
				// Append object to array
				phot.Comments = append(phot.Comments, commentObject)
			}
		}
	}
	for _, v := range photos {
		logrus.Infof("photoCommentsIdentifiers %v", v.ID)
	}
	return photos
}

func appendCommentCount(cnf config.Config, photoCommentCountIdentifiers []*sharedModels.CommentCountRequest, photos []*models.Photo) []*models.Photo {
	count := getCommentCount(cnf, photoCommentCountIdentifiers)
	for ind := 0; ind < len(photos); ind++ {
		// Get reference
		phot := photos[ind]
		for i := 0; i < len(count); i++ {
			countObject := count[i]
			if phot.ID == countObject.PhotoID {
				phot.CommentCount = countObject.Count
			}
		}
	}
	for _, v := range photos {
		logrus.Infof("photoCommentsIdentifiers2 %v", v.ID)
	}
	return photos
}

// appendUsernames triggers `GET usernames request` and appends result to []Photo
func appendUsernames(cnf config.Config, photoUsernamesIndentifiers []*sharedModels.GetUsernamesRequest, photos []*models.Photo) []*models.Photo {
	// Append the username to the photos
	usernames := getUsername(cnf, photoUsernamesIndentifiers)
	for index := 0; index < len(photos); index++ {
		photoObject := photos[index]
		for resultIndex := 0; resultIndex < len(usernames); resultIndex++ {
			userObject := usernames[resultIndex]
			if photoObject.UserID == userObject.ID {
				photoObject.Username = userObject.Username
			}
		} // end: for username/ids
	}
	return photos
}

// appendVotesCount triggers `GET votes request` and appends result to []Photo
func appendVotesCount(cnf config.Config, photoCountIdentifiers []*sharedModels.VoteCountRequest, photos []*models.Photo) []*models.Photo {
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
func appendUserVoted(cnf config.Config, photoVotedIdentifiers []*sharedModels.HasVotedRequest, photos []*models.Photo) []*models.Photo {
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

// Get the usernames from the ProfileService
func getUsername(cnf config.Config, input []*sharedModels.GetUsernamesRequest) []*sharedModels.GetUsernamesResponse {
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

		// Happy path
		data, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
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

// Get comments from CommentsService
func getComments(cnf config.Config, input []*sharedModels.CommentRequest) []*sharedModels.CommentResponse {
	type Req struct {
		Requests []*sharedModels.CommentRequest `json:"requests"`
	}
	body, _ := json.Marshal(Req{Requests: input})

	// Make url
	url := cnf.CommentServiceBaseurl + "ipc/getLast10"

	// Return object
	comments := make([]*sharedModels.CommentResponse, 0)

	// GET data and append to return object
	util.Request("GET", url, body, func(res *http.Response) {
		// Error handling
		if res.StatusCode < 200 || res.StatusCode > 299 {
			logrus.Errorf("ERROR with code %v.", res.Status)
			data, _ := ioutil.ReadAll(res.Body)
			logrus.Error(string(data))
			return
		}

		// Happy path
		data, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err == nil {
			jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				comment := &sharedModels.CommentResponse{}
				json.Unmarshal(value, comment)
				comments = append(comments, comment)
			}, "comments")
		}
	})
	return comments
}

// Get comment count from CommentsService
func getCommentCount(cnf config.Config, input []*sharedModels.CommentCountRequest) []*sharedModels.CommentCountResponse {
	type Req struct {
		Requests []*sharedModels.CommentCountRequest `json:"requests"`
	}
	body, _ := json.Marshal(Req{Requests: input})

	// Make url
	url := cnf.CommentServiceBaseurl + "ipc/getCount"

	// Return object
	comments := make([]*sharedModels.CommentCountResponse, 0)

	// GET data and append to return object
	util.Request("GET", url, body, func(res *http.Response) {
		// Error handling
		if res.StatusCode < 200 || res.StatusCode > 299 {
			logrus.Errorf("ERROR with code %v.", res.Status)
			data, _ := ioutil.ReadAll(res.Body)
			logrus.Error(string(data))
			return
		}

		// Happy path
		data, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err == nil {
			jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				comment := &sharedModels.CommentCountResponse{}
				json.Unmarshal(value, comment)
				comments = append(comments, comment)
			}, "result")
		}
	})
	return comments
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

	util.Request("GET", url, body, func(res *http.Response) {
		// Error handling
		if res.StatusCode < 200 || res.StatusCode > 299 {
			logrus.Errorf("ERROR with code %v.", res.Status)
			data, _ := ioutil.ReadAll(res.Body)
			logrus.Error(string(data))
			return
		}

		// Happy path
		data, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err == nil {
			jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				vote := &sharedModels.VoteCountResponse{}
				json.Unmarshal(value, vote)
				votes = append(votes, vote)
			}, "results")
		}
	})
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

	util.Request("GET", url, body, func(res *http.Response) {
		// Error handling
		if res.StatusCode < 200 || res.StatusCode > 299 {
			logrus.Errorf("ERROR with code %v.", res.Status)
			data, _ := ioutil.ReadAll(res.Body)
			logrus.Error(string(data))
			return
		}

		// Happy path
		data, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err == nil {
			jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				vote := &sharedModels.HasVotedResponse{}
				json.Unmarshal(value, vote)
				hasVoted = append(hasVoted, vote)
			}, "results")
		}
	})
	return hasVoted
}

func getUserIDFromRequest(cnf config.Config, req *http.Request) (int, error) {
	var queryToken = req.URL.Query().Get("token")

	if len(queryToken) < 1 {
		queryToken = req.Header.Get("token")
	}

	if len(queryToken) < 1 {
		return 0, errors.New("No JWT available")
	}

	tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(cnf.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims := tok.Claims.(jwt.MapClaims)
	var ID = claims["sub"].(float64) // gets the ID

	return int(ID), nil
}
