package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/config"

	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

// MariaDB struct for holding all methods related to the database
type MariaDB struct {
}

// mariaDBInstance is a private var. used as singleton
var mariaDBInstance *MariaDB

// InitMariaDB returns the instance of MariaDB
func InitMariaDB() *MariaDB {
	if mariaDBInstance == nil {
		mariaDBInstance = &MariaDB{}
	}
	return mariaDBInstance
}

// OpenConnection opens the connection to the database
func OpenConnection() (*sql.DB, error) {

	cnf := config.LoadConfig()

	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)

	log.Printf("Connect to : %v\n", dsn)
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

// CloseConnection closes the connection to the database
func CloseConnection(db *sql.DB) {
	db.Close()
}

// GetUserByUsername return the models.User object based on the username
func (mariaDB MariaDB) GetUserByUsername(username string) (models.User, error) {
	// Open the connection and close it with defer
	db, err := OpenConnection()
	if err != nil {
		return models.User{}, err
	}
	defer CloseConnection(db)

	// Query the database
	query := "SELECT id, username, createdAt, password, email FROM users WHERE username = '" + username + "'"
	rows, err := db.Query(query)
	if err != nil {
		return models.User{}, err
	}

	// Get first (and only) row
	if rows.Next() {
		var id int
		var username string
		var createdAt string
		var password string
		var email string
		err = rows.Scan(&id, &username, &createdAt, &password, &email)
		if err != nil {
			return models.User{}, err
		}

		return models.User{ID: id, Username: username, CreatedAt: util.TimeHelper(createdAt), Password: password, Email: email}, nil
	}
	return models.User{}, ErrUserNotFound
}

// ErrUserNotFound error if user does not exist in database
var ErrUserNotFound = errors.New("User does not exist")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
