package db

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
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

	cnf := config.LoadConfig()

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

func (mariaDB MariaDB) Create(comment *models.CommentCreate) (*sharedModels.CommentResponse, error) {
	db, err := OpenConnection()
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}
	defer CloseConnection(db)

	stmt, err := db.Prepare("INSERT INTO comments(user_id, photo_id, comment) VALUES(?,?,?)")
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}

	res, err := stmt.Exec(comment.UserID, comment.PhotoID, comment.Comment)
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}

	insertedID, err := res.LastInsertId()
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}

	c, err := mariaDB.GetCommentByID(int(insertedID))
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}
	return c, nil
}

func (mariaDB MariaDB) GetCommentByID(id int) (*sharedModels.CommentResponse, error) {
	db, err := OpenConnection()
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}
	defer CloseConnection(db)

	query := fmt.Sprintf("SELECT id, user_id, photo_id, comment, createdAt FROM comments WHERE id = %v", id)
	rows, err := db.Query(query)
	if err != nil {
		return &sharedModels.CommentResponse{}, err
	}
	if rows.Next() {
		comment := &sharedModels.CommentResponse{}
		var createdAt string

		rows.Scan(&comment.ID, &comment.UserID, &comment.PhotoID, &comment.Comment, &createdAt)

		comment.CreatedAt = util.TimeHelper(createdAt)

		return comment, nil
	}
	return nil, ErrCommentNotFound
}

func (mariaDB MariaDB) GetComments(photoID, offset, nrOfRows int) ([]*sharedModels.CommentResponse, error) {
	query := fmt.Sprintf("SELECT id, user_id, photo_id, comment, createdAt FROM comments WHERE photo_id=%v ORDER BY createdAt DESC LIMIT %v, %v", photoID, offset, nrOfRows)

	// Database actions
	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	responses := make([]*sharedModels.CommentResponse, 0)
	for rows.Next() {
		obj := &sharedModels.CommentResponse{}
		var createdAt string
		rows.Scan(&obj.ID, &obj.UserID, &obj.PhotoID, &obj.Comment, &createdAt)
		obj.CreatedAt = util.TimeHelper(createdAt)
		responses = append(responses, obj)
	}
	return responses, nil
}

func (mariaDB MariaDB) GetCommentCount(items []*sharedModels.CommentCountRequest) ([]*sharedModels.CommentCountResponse, error) {

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

	// Database actions
	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

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

func (mariaDB MariaDB) GetLastTenComments(items []*sharedModels.CommentRequest) ([]*sharedModels.CommentResponse, error) {
	if len(items) < 1 {
		return nil, nil
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

	// Database actions
	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	responses := make([]*sharedModels.CommentResponse, 0)
	for rows.Next() {
		obj := &sharedModels.CommentResponse{}
		var createdAt string
		rows.Scan(&obj.ID, &obj.UserID, &obj.PhotoID, &obj.Comment, &createdAt)
		obj.CreatedAt = util.TimeHelper(createdAt)
		responses = append(responses, obj)
	}
	return responses, nil
}

// ErrCommentNotFound error if comment does not exist in database
var ErrCommentNotFound = errors.New("Comment does not exist")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
