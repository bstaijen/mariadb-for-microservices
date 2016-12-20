package db

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"strconv"

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
func OpenConnection() *sql.DB {

	cnf := config.LoadConfig()

	username := cnf.DBUsername
	password := cnf.DBPassword
	host := cnf.DBHost
	port := cnf.DBPort
	database := cnf.Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)

	fmt.Printf("Connect to : %v\n", dsn)
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

func (mariaDB MariaDB) GetUserByID(ID int) (models.User, error) {
	db := OpenConnection()
	defer CloseConnection(db)

	query := "SELECT id, username, createdAt, password, email FROM users WHERE id = " + strconv.Itoa(ID)

	rows, err := db.Query(query)
	util.PanicIfError(err)

	if rows.Next() {
		var id int
		var username string
		var createdAt string
		var password string
		var email string
		err = rows.Scan(&id, &username, &createdAt, &password, &email)
		util.PanicIfError(err)

		return models.User{ID: id, Username: username, CreatedAt: util.TimeHelper(createdAt), Password: password, Email: email}, nil
	}
	return models.User{}, ErrUserNotFound
}

func (mariaDB MariaDB) GetUserByUsername(username string) (models.User, error) {
	db := OpenConnection()
	defer CloseConnection(db)

	query := "SELECT id, username, createdAt, password, email FROM users WHERE username = '" + username + "'"

	rows, err := db.Query(query)
	util.PanicIfError(err)

	if rows.Next() {
		var id int
		var username string
		var createdAt string
		var password string
		var email string
		err = rows.Scan(&id, &username, &createdAt, &password, &email)
		util.PanicIfError(err)

		return models.User{ID: id, Username: username, CreatedAt: util.TimeHelper(createdAt), Password: password, Email: email}, nil
	}
	return models.User{}, ErrUserNotFound
}

func (mariaDB MariaDB) CreateUser(user *models.User) (int, error) {
	db := OpenConnection()
	defer CloseConnection(db)

	// check unique username
	query := "SELECT * FROM users WHERE username = '" + user.Username + "'"
	rows, err := db.Query(query)
	util.PanicIfError(err)
	if rows.Next() {
		return 0, ErrUsernameIsNotUnique
	}

	// check unique email
	query = "SELECT * FROM users WHERE email = '" + user.Email + "'"
	rows, err = db.Query(query)
	util.PanicIfError(err)
	if rows.Next() {
		return 0, ErrEmailIsNotUnique
	}

	//Insert
	stmt, err := db.Prepare("INSERT INTO users(username, email, password) VALUES(?,?, ?)")
	util.PanicIfError(err)

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	res, err := stmt.Exec(user.Username, user.Email, string(hash))
	util.PanicIfError(err)

	id, err := res.LastInsertId()
	util.PanicIfError(err)

	return int(id), nil
}

func (mariaDB MariaDB) UpdateUser(user *models.User) int {
	db := OpenConnection()
	defer CloseConnection(db)

	stmt, err := db.Prepare("UPDATE users SET username = ?, email = ? WHERE id = ?")
	util.PanicIfError(err)

	_, err = stmt.Exec(user.Username, user.Email, user.ID)
	util.PanicIfError(err)

	return user.ID
}

func (mariaDB MariaDB) DeleteUser(user *models.User) int {
	db := OpenConnection()
	defer CloseConnection(db)

	if user.ID > 0 {

		stmt, err := db.Prepare("DELETE from users WHERE id = ?")
		util.PanicIfError(err)

		res, err := stmt.Exec(user.ID)
		util.PanicIfError(err)

		rowsAffected, err := res.RowsAffected()
		util.PanicIfError(err)

		return int(rowsAffected)
	}
	return -1

}

func (mariaDB MariaDB) GetUsers() []models.User {
	db := OpenConnection()
	defer CloseConnection(db)

	rows, err := db.Query("SELECT id, username, email, createdAt FROM users")
	util.PanicIfError(err)
	persons := make([]models.User, 0)

	for rows.Next() {
		var id int
		var username string
		var email string
		var createdAt string
		err = rows.Scan(&id, &username, &email, &createdAt)
		util.PanicIfError(err)

		persons = append(persons, models.User{ID: id, Username: username, CreatedAt: util.TimeHelper(createdAt)})
	}
	return persons
}

func (mariaDB MariaDB) GetUsernames(identifiers []*sharedModels.GetUsernamesRequest) []*sharedModels.GetUsernamesResponse {

	if len(identifiers) < 1 {
		return make([]*sharedModels.GetUsernamesResponse, 0)
	}

	query := inQueryBuilder(identifiers)

	db := OpenConnection()
	defer CloseConnection(db)

	rows, err := db.Query(query)
	util.PanicIfError(err)
	persons := make([]*sharedModels.GetUsernamesResponse, 0)

	for rows.Next() {
		var id int
		var username string
		err = rows.Scan(&id, &username)
		util.PanicIfError(err)

		persons = append(persons, &sharedModels.GetUsernamesResponse{ID: id, Username: username})
	}
	return persons
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
var ErrUserNotFound = errors.New("User does not exist")
