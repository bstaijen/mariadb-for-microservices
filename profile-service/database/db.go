package db

import (
	"database/sql"
	"errors"
	"fmt"

	"strconv"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

// MariaDB
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

func (mariaDB MariaDB) GetUserByID(ID int) (models.User, error) {
	db, err := OpenConnection()
	if err != nil {
		return models.User{}, err
	}
	defer CloseConnection(db)

	query := "SELECT id, username, createdAt, password, email FROM users WHERE id = " + strconv.Itoa(ID)

	rows, err := db.Query(query)
	if err != nil {
		return models.User{}, err
	}

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

func (mariaDB MariaDB) GetUserByUsername(username string) (models.User, error) {
	db, err := OpenConnection()
	if err != nil {
		return models.User{}, err
	}
	defer CloseConnection(db)

	query := "SELECT id, username, createdAt, password, email FROM users WHERE username = '" + username + "'"

	rows, err := db.Query(query)
	if err != nil {
		return models.User{}, err
	}

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

func (mariaDB MariaDB) CreateUser(user *models.User) (int, error) {
	db, err := OpenConnection()
	if err != nil {
		return 0, err
	}
	defer CloseConnection(db)

	// check unique username
	query := "SELECT * FROM users WHERE username = '" + user.Username + "'"
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		return 0, ErrUsernameIsNotUnique
	}

	// check unique email
	query = "SELECT * FROM users WHERE email = '" + user.Email + "'"
	rows, err = db.Query(query)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		return 0, ErrEmailIsNotUnique
	}

	//Insert
	stmt, err := db.Prepare("INSERT INTO users(username, email, password) VALUES(?,?, ?)")
	if err != nil {
		return 0, err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	res, err := stmt.Exec(user.Username, user.Email, string(hash))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (mariaDB MariaDB) UpdateUser(user *models.User) (int, error) {
	db, err := OpenConnection()
	if err != nil {
		return 0, err
	}
	defer CloseConnection(db)

	stmt, err := db.Prepare("UPDATE users SET username = ?, email = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}

	_, err = stmt.Exec(user.Username, user.Email, user.ID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (mariaDB MariaDB) DeleteUser(user *models.User) (int, error) {
	db, err := OpenConnection()
	if err != nil {
		return 0, err
	}
	defer CloseConnection(db)

	if user.ID > 0 {

		stmt, err := db.Prepare("DELETE from users WHERE id = ?")
		if err != nil {
			return 0, err
		}

		res, err := stmt.Exec(user.ID)
		if err != nil {
			return 0, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		return int(rowsAffected), nil
	}
	return 0, errors.New("User ID is empty")

}

func (mariaDB MariaDB) GetUsers() ([]models.User, error) {
	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

	rows, err := db.Query("SELECT id, username, email, createdAt FROM users")
	if err != nil {
		return nil, err
	}
	persons := make([]models.User, 0)

	for rows.Next() {
		var id int
		var username string
		var email string
		var createdAt string
		err = rows.Scan(&id, &username, &email, &createdAt)
		if err != nil {
			return nil, err
		}
		persons = append(persons, models.User{ID: id, Username: username, CreatedAt: util.TimeHelper(createdAt)})
	}
	return persons, nil
}

func (mariaDB MariaDB) GetUsernames(identifiers []*sharedModels.GetUsernamesRequest) ([]*sharedModels.GetUsernamesResponse, error) {

	if len(identifiers) < 1 {
		return make([]*sharedModels.GetUsernamesResponse, 0), nil
	}

	query := inQueryBuilder(identifiers)

	db, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(db)

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
