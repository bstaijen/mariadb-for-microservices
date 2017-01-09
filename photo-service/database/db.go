package db

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
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

func OpenConnection() (*sql.DB, error) {

	cnf := config.LoadConfig()

	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)

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

func (mariaDB MariaDB) InsertPhoto(photo *models.CreatePhoto) error {

	db, err := OpenConnection()
	if err != nil {
		return err
	}
	defer CloseConnection(db)

	//Insert
	stmt, err := db.Prepare("INSERT INTO photos(user_id, filename, title, contentType, photo) VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(photo.UserID, photo.Filename, photo.Title, photo.ContentType, photo.Image)
	if err != nil {
		return err
	}
	return nil
}

func (mariaDB MariaDB) ListImagesByUserID(id int) ([]*models.Photo, error) {
	return selectQuery(fmt.Sprintf("SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos WHERE user_id=%v", id))
}

func (mariaDB MariaDB) ListIncoming(offset int, nrOfRows int) ([]*models.Photo, error) {
	return selectQuery(fmt.Sprintf("SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos ORDER BY createdAt DESC LIMIT %v, %v", offset, nrOfRows))
}

func (mariaDB MariaDB) GetPhotoByFilename(filename string) (*models.Photo, error) {
	photos, err := selectQuery(fmt.Sprintf("SELECT id, user_id, filename, title, createdAt, contentType, photo FROM photos WHERE filename = '%v'", filename))
	if len(photos) > 0 {
		return photos[0], err
	}
	return nil, err
}

func (mariaDB MariaDB) GetPhotoById(id int) (*models.Photo, error) {
	photos, err := selectQuery(fmt.Sprintf("SELECT id, user_id, filename, title, createdAt, contentType, photo  FROM photos WHERE id = %v", id))
	if len(photos) > 0 {
		return photos[0], err
	}
	return nil, nil
}

func selectQuery(query string) ([]*models.Photo, error) {

	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	photos := []*models.Photo{}
	for rows.Next() {
		photoObject := &models.Photo{}

		var createdAt string

		err = rows.Scan(&photoObject.ID, &photoObject.UserID, &photoObject.Filename, &photoObject.Title, &createdAt, &photoObject.ContentType, &photoObject.Image)
		if err != nil {
			return nil, err
		}

		photoObject.CreatedAt = util.TimeHelper(createdAt)

		photos = append(photos, photoObject)
	}
	return photos, nil
}

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
