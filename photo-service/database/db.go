package db

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
)

// OpenConnection opens the connection to the database
func OpenConnection() (*sql.DB, error) {

	cnf := config.LoadConfig()

	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", username, password, host, port, database)

	log.Debugf("Connect to : %v\n", dsn)
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

func InsertPhoto(db *sql.DB, photo *models.CreatePhoto) error {
	//Insert
	_, err := db.Exec("INSERT INTO photos(user_id, filename, title, contentType, photo) VALUES(?,?,?,?,?)", photo.UserID, photo.Filename, photo.Title, photo.ContentType, photo.Image)
	if err != nil {
		return err
	}
	return nil
}

func ListImagesByUserID(db *sql.DB, id int) ([]*models.Photo, error) {
	return selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos WHERE user_id=?", id)
}

func ListIncoming(db *sql.DB, offset int, nrOfRows int) ([]*models.Photo, error) {
	return selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos ORDER BY createdAt DESC LIMIT ?, ?", offset, nrOfRows)
}

func GetPhotoByFilename(db *sql.DB, filename string) (*models.Photo, error) {
	photos, err := selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos WHERE filename = ?", filename)
	if len(photos) > 0 {
		return photos[0], err
	}
	return nil, err
}

func GetPhotoById(db *sql.DB, id int) (*models.Photo, error) {
	photos, err := selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo  FROM photos WHERE id = ?", id)
	if len(photos) > 0 {
		return photos[0], err
	}
	return nil, nil
}

// A parameter type prefixed with three dots (...) is called a variadic parameter.
func selectQuery(db *sql.DB, query string, args ...interface{}) ([]*models.Photo, error) {

	log.Infof(query, args)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	photos := []*models.Photo{}
	for rows.Next() {
		photoObject := &models.Photo{}

		err = rows.Scan(&photoObject.ID, &photoObject.UserID, &photoObject.Filename, &photoObject.Title, &photoObject.CreatedAt, &photoObject.ContentType, &photoObject.Image)
		if err != nil {
			return nil, err
		}

		photos = append(photos, photoObject)
	}
	return photos, nil
}

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
