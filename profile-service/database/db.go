package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// OpenConnection method
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

func GetUserByID(db *sql.DB, ID int) (models.User, error) {
	rows, err := db.Query("SELECT id, username, createdAt, password, email FROM users WHERE id = ?", ID)
	if err != nil {
		return models.User{}, err
	}

	if rows.Next() {
		var id int
		var username string
		var createdAt time.Time
		var password string
		var email string
		err = rows.Scan(&id, &username, &createdAt, &password, &email)
		if err != nil {
			return models.User{}, err
		}

		return models.User{ID: id, Username: username, CreatedAt: createdAt, Password: password, Email: email}, nil
	}
	return models.User{}, ErrUserNotFound
}

func CreateUser(db *sql.DB, user *models.User) (int, error) {
	// check unique username
	query := "SELECT * FROM users WHERE username = ?"
	rows, err := db.Query(query, user.Username)
	if err != nil {
		log.Errorf("Error executing query %v ", query)
		log.Error(err)
		return 0, err
	}
	if rows.Next() {
		return 0, ErrUsernameIsNotUnique
	}

	// check unique email
	query = "SELECT * FROM users WHERE email = ?"
	rows, err = db.Query(query, user.Email)
	if err != nil {
		log.Errorf("Error executing query %v ", query)
		log.Error(err)
		return 0, err
	}
	if rows.Next() {
		return 0, ErrEmailIsNotUnique
	}

	// Insert
	res, err := db.Exec("INSERT INTO users (username, email, password) VALUES(?, ?, ?)", user.Username, user.Email, user.Hash)
	if err != nil {
		log.Errorf("Error inserting")
		log.Error(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return int(id), nil
}

func UpdateUser(db *sql.DB, user *models.User) (int, error) {
	_, err := db.Exec("UPDATE users SET username = ?, email = ? WHERE id = ?", user.Username, user.Email, user.ID)
	if err != nil {
		log.Errorf("Error inserting")
		log.Error(err)
		return 0, err
	}
	return user.ID, nil
}

func DeleteUser(db *sql.DB, user *models.User) (int, error) {
	if user.ID > 0 {
		res, err := db.Exec("DELETE from users WHERE id = ?", user.ID)
		if err != nil {
			log.Errorf("Error inserting")
			log.Error(err)
			return 0, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Errorf("Error inserting")
			log.Error(err)
			return 0, err
		}
		return int(rowsAffected), nil
	}
	return 0, errors.New("User ID is empty")

}

func GetUsers(db *sql.DB) ([]models.User, error) {

	rows, err := db.Query("SELECT id, username, email, createdAt FROM users")
	if err != nil {
		return nil, err
	}
	persons := make([]models.User, 0)

	for rows.Next() {
		var id int
		var username string
		var email string
		var createdAt time.Time
		err = rows.Scan(&id, &username, &email, &createdAt)
		if err != nil {
			return nil, err
		}
		persons = append(persons, models.User{ID: id, Username: username, CreatedAt: createdAt})
	}
	return persons, nil
}

func GetUsernames(db *sql.DB, identifiers []*sharedModels.GetUsernamesRequest) ([]*sharedModels.GetUsernamesResponse, error) {

	if len(identifiers) < 1 {
		return make([]*sharedModels.GetUsernamesResponse, 0), nil
	}

	query := inQueryBuilder(identifiers)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	persons := make([]*sharedModels.GetUsernamesResponse, 0)

	for rows.Next() {
		var id int
		var username string
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, err
		}

		persons = append(persons, &sharedModels.GetUsernamesResponse{ID: id, Username: username})
	}
	return persons, nil
}

// Query builder for contructing an IN-condition
func inQueryBuilder(identifiers []*sharedModels.GetUsernamesRequest) string {
	if len(identifiers) < 1 {
		return ""
	}

	query := "SELECT id, username FROM users WHERE id IN"
	query += "("

	for i := 0; i < len(identifiers); i++ {
		if i+1 < len(identifiers) {
			// NOT LAST
			query += strconv.Itoa(identifiers[i].ID) + ","
		} else {
			//LAST
			query += strconv.Itoa(identifiers[i].ID)
		}
	}

	query += ")"
	return query
}

var ErrEmailIsNotUnique = errors.New("Email must be unique")
var ErrUsernameIsNotUnique = errors.New("Username must be unique")

// ErrUserNotFound error if user does not exist in database
var ErrUserNotFound = errors.New("User does not exist")

// ErrCanNotConnectWithDatabase error if database is unreachable
var ErrCanNotConnectWithDatabase = errors.New("Can not connect with database")
