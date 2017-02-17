package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/vote-service/config"
)

// OpenConnection opens the connection to the database
func OpenConnection(cnf config.Config) (*sql.DB, error) {
	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)

	log.Debugf("Connect to : %v", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, ErrCanNotConnectWithDatabase
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, ErrCanNotConnectWithDatabase
	}
	return db, nil
}

// CloseConnection method
func CloseConnection(db *sql.DB) {
	db.Close()
}

// Create a vote in the database
func Create(db *sql.DB, vote *sharedModels.VoteCreateRequest) error {
	_, err := db.Exec("DELETE FROM votes WHERE user_id=? AND photo_id=?", vote.UserID, vote.PhotoID)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO votes(user_id, photo_id, upvote, downvote) VALUES(?,?,?,?)", vote.UserID, vote.PhotoID, vote.Upvote, vote.Downvote)
	if err != nil {
		return err
	}
	return nil
}

// VoteCount calculates how many votes each photos has.
func VoteCount(db *sql.DB, items []*sharedModels.VoteCountRequest) ([]*sharedModels.VoteCountResponse, error) {
	if len(items) < 1 {
		return make([]*sharedModels.VoteCountResponse, 0), nil
	}

	// QUERY BUILDER
	query := "SELECT photo_id, sum(upvote), sum(downvote) FROM votes WHERE photo_id IN"
	query += "("

	for i := 0; i < len(items); i++ {
		if i+1 < len(items) {
			// NOT LAST
			query += strconv.Itoa(items[i].PhotoID) + ","
		} else {
			//LAST
			query += strconv.Itoa(items[i].PhotoID)
		}
	}

	query += ") GROUP BY photo_id"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	photoCountResponses := make([]*sharedModels.VoteCountResponse, 0)

	for rows.Next() {
		obj := &sharedModels.VoteCountResponse{}
		err = rows.Scan(&obj.PhotoID, &obj.UpVoteCount, &obj.DownVoteCount)
		if err != nil {
			return nil, err
		}
		photoCountResponses = append(photoCountResponses, obj)
	}
	return photoCountResponses, nil
}

// HasVoted is a method which calculates whether or not a user has voted on a photo. Method accepts a list of photos and returns the result for each photo in the list.
func HasVoted(db *sql.DB, items []*sharedModels.HasVotedRequest) ([]*sharedModels.HasVotedResponse, error) {
	if len(items) < 1 {
		return make([]*sharedModels.HasVotedResponse, 0), nil
	}

	// QUERY BUILDER
	query := "SELECT user_id, photo_id, upvote, downvote FROM votes WHERE "
	query += ""

	for i := 0; i < len(items); i++ {
		if i+1 < len(items) {
			// NOT LAST
			query += fmt.Sprintf("(photo_id = %v AND user_id = %v) OR ", strconv.Itoa(items[i].PhotoID), strconv.Itoa(items[i].UserID))
		} else {
			//LAST
			query += fmt.Sprintf("(photo_id = %v AND user_id = %v)", strconv.Itoa(items[i].PhotoID), strconv.Itoa(items[i].UserID))
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	photoVotedResponses := make([]*sharedModels.HasVotedResponse, 0)

	for rows.Next() {
		obj := &sharedModels.HasVotedResponse{}
		// obj.Voted = true

		err = rows.Scan(&obj.UserID, &obj.PhotoID, &obj.Upvote, &obj.Downvote)
		if err != nil {
			return nil, err
		}

		photoVotedResponses = append(photoVotedResponses, obj)
	}
	return photoVotedResponses, nil
}

// GetTopRatedTimeline returns an array of top rated photos. The array contains a list of ID's. Offset and nrOfRows can be used for pagination.
func GetTopRatedTimeline(db *sql.DB, offset int, nrOfRows int) ([]*sharedModels.TopRatedPhotoResponse, error) {
	rows, err := db.Query("SELECT photo_id as photoID, sum(upvote) AS totalUpvote, sum(downvote) AS totalDownvote, sum(upvote) - sum(downvote) as difference FROM votes GROUP BY photo_id ORDER BY difference DESC LIMIT ?, ?", offset, nrOfRows)
	if err != nil {
		return nil, err
	}

	photos := make([]*sharedModels.TopRatedPhotoResponse, 0)

	for rows.Next() {
		var photoID int
		var totalUpvote int
		var totalDownvote int
		var difference int

		err = rows.Scan(&photoID, &totalUpvote, &totalDownvote, &difference)
		if err != nil {
			return nil, err
		}
		photos = append(photos, &sharedModels.TopRatedPhotoResponse{
			PhotoID: photoID,
		})
	}
	return photos, nil
}

// GetHotTimeline returns an array of photos ordered by on which is most 'hot' meaning which has been voted on the most for the CURRENT_DAY
func GetHotTimeline(db *sql.DB, offset int, nrOfRows int) ([]*sharedModels.TopRatedPhotoResponse, error) {
	rows, err := db.Query("SELECT photo_id as photoID, sum(upvote) AS totalUpvote, sum(downvote) AS totalDownvote, sum(upvote) - sum(downvote) AS difference FROM votes WHERE createdAt > DATE_SUB(now(), INTERVAL 1 DAY) GROUP BY photo_id ORDER BY difference DESC LIMIT ?, ?", offset, nrOfRows)
	if err != nil {
		return nil, err
	}

	photos := make([]*sharedModels.TopRatedPhotoResponse, 0)

	for rows.Next() {
		var photoID int
		var totalUpvote int
		var totalDownvote int
		var difference int

		err = rows.Scan(&photoID, &totalUpvote, &totalDownvote, &difference)
		if err != nil {
			return nil, err
		}

		photos = append(photos, &sharedModels.TopRatedPhotoResponse{
			PhotoID: photoID,
		})
	}
	return photos, nil
}

// GetVotesFromUser returns the votes the user has placed on photos. Order by last created.
func GetVotesFromUser(db *sql.DB, userID int, offset int, nrOfRows int) ([]*sharedModels.TopRatedPhotoResponse, error) {
	rows, err := db.Query("SELECT photo_id FROM votes WHERE user_id = ? ORDER BY createdAt DESC LIMIT ?, ?", userID, offset, nrOfRows)
	if err != nil {
		return nil, err
	}
	photos := make([]*sharedModels.TopRatedPhotoResponse, 0)
	for rows.Next() {
		var photoID int
		err = rows.Scan(&photoID)
		if err != nil {
			return nil, err
		}
		photos = append(photos, &sharedModels.TopRatedPhotoResponse{
			PhotoID: photoID,
		})
	}
	return photos, nil
}

// ErrUserNotFound error if user does not exist in database
var ErrUserNotFound = errors.New("User does not exist")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
