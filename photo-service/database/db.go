package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// OpenConnection opens the connection to the database
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
		return nil, errCanNotConnectWithDatabase
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, errCanNotConnectWithDatabase
	}
	return db, nil
}

// CloseConnection method
func CloseConnection(db *sql.DB) {
	db.Close()
}

// InsertPhoto : inserts a photo in the database
func InsertPhoto(db *sql.DB, photo *models.CreatePhoto) error {
	//Insert
	_, err := db.Exec("INSERT INTO photos(user_id, filename, title, contentType, photo) VALUES(?,?,?,?,?)", photo.UserID, photo.Filename, photo.Title, photo.ContentType, photo.Image)
	if err != nil {
		return err
	}
	return nil
}

// ListImagesByUserID returns a list of photo's uploaded by the user.
func ListImagesByUserID(db *sql.DB, id int) ([]*models.Photo, error) {
	return selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos WHERE user_id=? ORDER BY createdAt DESC", id)
}

// ListIncoming returns a list of photos ordered by last inserted
func ListIncoming(db *sql.DB, offset int, nrOfRows int) ([]*models.Photo, error) {
	return selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos ORDER BY createdAt DESC LIMIT ?, ?", offset, nrOfRows)
}

// GetPhotoByFilename return a photo based on the filename
func GetPhotoByFilename(db *sql.DB, filename string) (*models.Photo, error) {
	photos, err := selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos WHERE filename = ?", filename)
	if len(photos) > 0 {
		return photos[0], err
	}
	return nil, err 
}

// GetPhotoById returns a photo indexed by id
func GetPhotoById(db *sql.DB, id int) (*models.Photo, error) {
	photos, err := selectQuery(db, "SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos WHERE id = ?", id)

	log.Info(photos)

	if len(photos) > 0 {
		return photos[0], err
	} else if len(photos) == 0 {
		return nil, errors.New("photo not found")
	}
	return nil, err
}

func GetPhotos(db *sql.DB, items []*sharedModels.PhotoRequest) ([]*sharedModels.PhotoResponse, error) {
	if len(items) < 1 {
		return make([]*sharedModels.PhotoResponse, 0), nil
	}

	// QUERY BUILDER
	query := "SELECT id, user_id, filename, title, createdAt FROM photos WHERE id IN"
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

	query += ")"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	photos := []*sharedModels.PhotoResponse{}
	for rows.Next() {
		photoObject := &sharedModels.PhotoResponse{}

		err = rows.Scan(&photoObject.ID, &photoObject.UserID, &photoObject.Filename, &photoObject.Title, &photoObject.CreatedAt)
		if err != nil {
			return nil, err
		}

		photos = append(photos, photoObject)
	}
	return photos, nil
}

// DeletePhotoByID delete a photo in the database based on ID.
func DeletePhotoByID(db *sql.DB, photoID int) (int64, error) {
	res, err := db.Exec("DELETE FROM photos WHERE id = ?", photoID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// A parameter type prefixed with three dots (...) is called a variadic parameter.
func selectQuery(db *sql.DB, query string, args ...interface{}) ([]*models.Photo, error) {
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

// errCanNotConnectWithDatabase error if database is unreachable
var errCanNotConnectWithDatabase = errors.New("Can not connect with database")
