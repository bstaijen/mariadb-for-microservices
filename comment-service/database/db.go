package db

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// OpenConnection method
func OpenConnection(cnf config.Config) (*sql.DB, error) {
	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", username, password, host, port, database)

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

// Create : saves a comment in the database
func Create(db *sql.DB, comment *models.CommentCreate) (*sharedModels.CommentResponse, error) {
	res, err := db.Exec("INSERT INTO comments(user_id, photo_id, comment) VALUES(?,?,?)", comment.UserID, comment.PhotoID, comment.Comment)
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}

	insertedID, err := res.LastInsertId()
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}

	c, err := GetCommentByID(db, int(insertedID))
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}
	return c, nil
}

// GetCommentByID returns a comment from the database
func GetCommentByID(db *sql.DB, id int) (*sharedModels.CommentResponse, error) {
	rows, err := db.Query("SELECT id, user_id, photo_id, comment, createdAt FROM comments WHERE id = ?", id)
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}
	if rows.Next() {
		comment := &sharedModels.CommentResponse{}
		rows.Scan(&comment.ID, &comment.UserID, &comment.PhotoID, &comment.Comment, &comment.CreatedAt)
		return comment, nil
	}
	return nil, ErrCommentNotFound
}

// GetComments return an array of comments.
func GetComments(db *sql.DB, photoID, offset, nrOfRows int) ([]*sharedModels.CommentResponse, error) {
	rows, err := db.Query("SELECT id, user_id, photo_id, comment, createdAt FROM comments WHERE photo_id=? ORDER BY createdAt DESC LIMIT ?, ?", photoID, offset, nrOfRows)
	if err != nil {
		return nil, err
	}

	responses := make([]*sharedModels.CommentResponse, 0)
	for rows.Next() {
		obj := &sharedModels.CommentResponse{}
		rows.Scan(&obj.ID, &obj.UserID, &obj.PhotoID, &obj.Comment, &obj.CreatedAt)
		responses = append(responses, obj)
	}
	return responses, nil
}

// GetCommentCount returns the number of comment counts
func GetCommentCount(db *sql.DB, items []*sharedModels.CommentCountRequest) ([]*sharedModels.CommentCountResponse, error) {

	if len(items) < 1 {
		return nil, nil
	}
	// Query builder
	query := ""
	for index := 0; index < len(items); index++ {
		if index+1 < len(items) {
			// NOT LAST
			query += fmt.Sprintf("(SELECT photo_id, COUNT(*) FROM comments WHERE photo_id =%v) UNION ALL ", items[index].PhotoID)
		} else {
			//LAST
			query += fmt.Sprintf("(SELECT photo_id, COUNT(*) FROM comments WHERE photo_id =%v)", items[index].PhotoID)
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	responses := make([]*sharedModels.CommentCountResponse, 0)
	for rows.Next() {
		obj := &sharedModels.CommentCountResponse{}
		rows.Scan(&obj.PhotoID, &obj.Count)
		responses = append(responses, obj)
	}
	return responses, nil
}

// GetLastTenComments return the 10 comments for each comment in the [] parameter.
func GetLastTenComments(db *sql.DB, items []*sharedModels.CommentRequest) ([]*sharedModels.CommentResponse, error) {
	responses := make([]*sharedModels.CommentResponse, 0)

	if len(items) < 1 {
		return responses, nil
	}
	// Query builder
	query := ""
	for index := 0; index < len(items); index++ {
		if index+1 < len(items) {
			// NOT LAST
			query += fmt.Sprintf("(SELECT id, user_id, photo_id, comment, createdAt FROM comments WHERE photo_id=%v ORDER BY createdAt DESC LIMIT 10) UNION ALL ", items[index].PhotoID)
		} else {
			//LAST
			query += fmt.Sprintf("(SELECT id, user_id, photo_id, comment, createdAt FROM comments WHERE photo_id=%v ORDER BY createdAt DESC LIMIT 10)", items[index].PhotoID)
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		obj := &sharedModels.CommentResponse{}
		rows.Scan(&obj.ID, &obj.UserID, &obj.PhotoID, &obj.Comment, &obj.CreatedAt)
		responses = append(responses, obj)
	}
	return responses, nil
}

// ErrCommentNotFound error if comment does not exist in database
var ErrCommentNotFound = errors.New("Comment does not exist")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
