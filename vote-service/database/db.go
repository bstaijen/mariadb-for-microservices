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

type MariaDB struct {
}

var mariaDBInstance *MariaDB = nil

func InitMariaDB() *MariaDB {
	if mariaDBInstance == nil {
		mariaDBInstance = &MariaDB{}
	}
	return mariaDBInstance
}

// OpenConnection method
func OpenConnection() (*sql.DB, error) {

	cnf := config.LoadConfig() // what if this breaks ?

	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)

	log.Debug("Connect to : %v\n", dsn)
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

func (mariaDB MariaDB) Create(vote *sharedModels.VoteCreateRequest) error {
	db, err := OpenConnection()
	if err != nil {
		return err
	}
	defer CloseConnection(db)

	stmt, err := db.Prepare("DELETE FROM votes WHERE user_id=? AND photo_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vote.UserID, vote.PhotoID)
	if err != nil {
		return err
	}
	stmt, err = db.Prepare("INSERT INTO votes(user_id, photo_id, upvote, downvote) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(vote.UserID, vote.PhotoID, vote.Upvote, vote.Downvote)
	if err != nil {
		return err
	}
	return nil
}

func (mariaDB MariaDB) VoteCount(items []*sharedModels.VoteCountRequest) ([]*sharedModels.VoteCountResponse, error) {
	if len(items) < 1 {
		return make([]*sharedModels.VoteCountResponse, 0), nil
	}

	// QUERY BUILDER
	query := "SELECT photo_id, sum(upvote), sum(downvote) FROM votes WHERE photo_id IN"
	query += "("

	for i := 0; i < len(items); i++ { // xx any oppportunities for sql injection here ?
		if i+1 < len(items) {
			// NOT LAST
			query += strconv.Itoa(items[i].PhotoID) + ","
		} else {
			//LAST
			query += strconv.Itoa(items[i].PhotoID)
		}
	}

	query += ") GROUP BY photo_id"

	// DO DATABASE THINGS
	db, err := OpenConnection() // xx  how often do we open / close databsases per request ?
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

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

func (mariaDB MariaDB) HasVoted(items []*sharedModels.HasVotedRequest) ([]*sharedModels.HasVotedResponse, error) {
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

	// DO DATABASE THINGS
	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

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

func (mariaDB MariaDB) GetTopRatedTimeline(offset int, nrOfRows int) ([]*sharedModels.TopRatedPhotoResponse, error) {
	query := fmt.Sprintf("select photo_id, sum(upvote) as total_upvote, sum(downvote) as total_downvote, sum(upvote) - sum(downvote) as difference from votes group by photo_id order by difference desc LIMIT %v, %v", offset, nrOfRows)

	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	photos := make([]*sharedModels.TopRatedPhotoResponse, 0)

	for rows.Next() {
		var photo_id int
		var total_upvote int
		var total_downvote int
		var difference int

		err = rows.Scan(&photo_id, &total_upvote, &total_downvote, &difference)
		if err != nil {
			return nil, err
		}
		photos = append(photos, &sharedModels.TopRatedPhotoResponse{
			PhotoID: photo_id,
		})
	}
	return photos, nil
}

func (mariaDB MariaDB) GetHotTimeline(offset int, nrOfRows int) ([]*sharedModels.TopRatedPhotoResponse, error) {
	query := fmt.Sprintf("select photo_id, sum(upvote) as total_upvote, sum(downvote) as total_downvote, sum(upvote) - sum(downvote) as difference from votes where createdAt > DATE_SUB(now(), INTERVAL 1 DAY) group by photo_id order by difference desc LIMIT %v, %v", offset, nrOfRows)
	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	photos := make([]*sharedModels.TopRatedPhotoResponse, 0)

	for rows.Next() {
		var photo_id int
		var total_upvote int
		var total_downvote int
		var difference int

		err = rows.Scan(&photo_id, &total_upvote, &total_downvote, &difference)
		if err != nil {
			return nil, err
		}

		photos = append(photos, &sharedModels.TopRatedPhotoResponse{
			PhotoID: photo_id,
		})
	}
	return photos, nil
}

// ErrUserNotFound error if user does not exist in database
var ErrUserNotFound = errors.New("User does not exist")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
