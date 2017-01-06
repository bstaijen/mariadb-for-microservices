package db

import (
	"database/sql"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
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

func (mariaDB MariaDB) Create(vote *sharedModels.VoteCreateRequest) error {
	db := OpenConnection()
	defer CloseConnection(db)

	stmt, err := db.Prepare("DELETE FROM votes WHERE user_id=? AND photo_id=?")
	util.PanicIfError(err)

	_, err = stmt.Exec(vote.UserID, vote.PhotoID)
	util.PanicIfError(err)

	stmt, err = db.Prepare("INSERT INTO votes(user_id, photo_id, upvote, downvote) VALUES(?,?,?,?)")
	util.PanicIfError(err)

	_, err = stmt.Exec(vote.UserID, vote.PhotoID, vote.Upvote, vote.Downvote)
	if err != nil {
		return err
	}
	return nil
}

func (mariaDB MariaDB) VoteCount(items []*sharedModels.VoteCountRequest) []*sharedModels.VoteCountResponse {
	if len(items) < 1 {
		return nil
	}

	// QUERY BUILDER
	// select photo_id, count(*) as count from votes where photo_id in (21,22) group by photo_id;
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
	db := OpenConnection() // xx  how often do we open / close databsases per request ?
	defer CloseConnection(db)

	rows, err := db.Query(query)
	util.PanicIfError(err)

	photoCountResponses := make([]*sharedModels.VoteCountResponse, 0)

	for rows.Next() {
		obj := &sharedModels.VoteCountResponse{}
		rows.Scan(&obj.PhotoID, &obj.UpVoteCount, &obj.DownVoteCount)
		photoCountResponses = append(photoCountResponses, obj)
	}
	return photoCountResponses
}

func (mariaDB MariaDB) HasVoted(items []*sharedModels.HasVotedRequest) []*sharedModels.HasVotedResponse {
	if len(items) < 1 {
		return nil
	}

	// QUERY BUILDER
	// select user_id, photo_id from votes where photo_id in (1,19,20,21,22) AND user_id = 2;
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

	query += ""

	// DO DATABASE THINGS
	db := OpenConnection()
	defer CloseConnection(db)

	rows, err := db.Query(query)
	util.PanicIfError(err)

	photoVotedResponses := make([]*sharedModels.HasVotedResponse, 0)

	for rows.Next() {
		obj := &sharedModels.HasVotedResponse{}
		// obj.Voted = true

		rows.Scan(&obj.UserID, &obj.PhotoID, &obj.Upvote, &obj.Downvote)

		photoVotedResponses = append(photoVotedResponses, obj)
	}
	return photoVotedResponses
}

func (mariaDB MariaDB) GetTopRatedTimeline() []*sharedModels.TopRatedPhotoResponse {
	query := "select photo_id, sum(upvote) as total_upvote, sum(downvote) as total_downvote, sum(upvote) - sum(downvote) as difference from votes group by photo_id order by difference desc limit 10"

	db := OpenConnection()
	defer CloseConnection(db)

	rows, err := db.Query(query)
	util.PanicIfError(err)

	photos := make([]*sharedModels.TopRatedPhotoResponse, 0)

	for rows.Next() {
		var photo_id int
		var total_upvote int
		var total_downvote int
		var difference int

		rows.Scan(&photo_id, &total_upvote, &total_downvote, &difference)

		photos = append(photos, &sharedModels.TopRatedPhotoResponse{
			PhotoID: photo_id,
		})
	}
	return photos
}

func (mariaDB MariaDB) GetHotTimeline() []*sharedModels.TopRatedPhotoResponse {
	query := "select photo_id, sum(upvote) as total_upvote, sum(downvote) as total_downvote, sum(upvote) - sum(downvote) as difference from votes where createdAt > DATE_SUB(now(), INTERVAL 1 DAY) group by photo_id order by difference desc limit 10"
	db := OpenConnection()
	defer CloseConnection(db)

	rows, err := db.Query(query)
	util.PanicIfError(err)

	photos := make([]*sharedModels.TopRatedPhotoResponse, 0)

	for rows.Next() {
		var photo_id int
		var total_upvote int
		var total_downvote int
		var difference int

		rows.Scan(&photo_id, &total_upvote, &total_downvote, &difference)

		photos = append(photos, &sharedModels.TopRatedPhotoResponse{
			PhotoID: photo_id,
		})
	}
	return photos
}

// OpenConnection method
func OpenConnection() *sql.DB {

	cnf := config.LoadConfig() // what if this breaks ?

	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)

	log.Debugf("Connect to : %v\n", dsn) // log vs printf
	db, err := sql.Open("mysql", dsn)
	util.PanicIfError(err)

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	util.PanicIfError(err)

	return db
}

// CloseConnection method
func CloseConnection(db *sql.DB) {
	db.Close()
}
