package controllers

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestHash struct{}

func (a TestHash) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

// Test creating an user.
func TestCreateUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expected rows
	rows := sqlmock.NewRows([]string{"count(*)"})

	// Expectation: check for unique username
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Username).WillReturnRows(rows)

	// Expectation: check for unique email
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Email).WillReturnRows(rows)

	// Expectation: insert into database
	mock.ExpectExec("INSERT INTO users").WithArgs(user.Username, user.Email, TestHash{}).WillReturnResult(sqlmock.NewResult(1, 1))

	// Expectation: get user by id
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "email"}).AddRow(user.ID, user.Username, timeNow, user.Email)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response is alright
	type Token struct {
		Token     string              `json:"token"`
		ExpiresOn string              `json:"expires_on"`
		User      models.UserResponse `json:"user"`
	}

	response := &Token{}
	err = decodeJSON(res.Body, response)
	if err != nil {
		t.Fatal(errors.New("Bad json"))
	}
	if response.User.ID < 1 {
		t.Errorf("Expected user ID greater than 0 but got %v", response.User.ID)
	}

	if response.User.Username != user.Username {
		t.Errorf("Expected username to be %v but got %v", user.Username, response.User.Username)
	}
	if response.User.Email != user.Email {
		t.Errorf("Expected username to be %v but got %v", user.Email, response.User.Email)
	}

	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}

// Test creating an user when a bad json string is provided. We expect an error message.
func TestBadJson(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer([]byte("{")))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Bad json\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Test creating a user without providing a username. We expect an error message.
func TestCreateUserWithoutUsername(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Username is too short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Test creating an user without providing a password. We expect an error message.
func TestCreateUserWithoutPassword(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}
	user.Email = "test@example.com"
	user.Username = "username"

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Password is to short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Test creating an user without providing an email address. We expect an error message.
func TestCreateUserWithoutEmail(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}
	user.Username = "username"

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := CreateUserHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Email address is to short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// Helper function to decode a json string to an interface
func decodeJSON(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		fmt.Printf("json decoder error occurred: %v \n", err.Error())
		return err
	}
	return nil
}

// Test deleting an user.
func TestDeleteUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.UserCreate{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"

	// Create JWT object with claims
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	secretKey := cnf.SecretKey
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Error(err)
		return
	}
	json, _ := json.Marshal(user)
	req, err := http.NewRequest("DELETE", "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE from users WHERE").WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	handler := DeleteUserHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}

// Test updating an user.
func TestUpdateUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := getTestUser()
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	req, err := http.NewRequest("PUT", "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE users SET").WithArgs(user.Username, user.Email, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	handler := UpdateUserHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}

// Test converting a json string to a list with ID's
func TestBodyToArrayWithIDs(t *testing.T) {
	mock := []byte(`{ "requests":[{"id":1} ,{"id":2},{"id":3}, {"id":4} ]}`)
	req, err := http.NewRequest("GET", "http://localhost/", bytes.NewBuffer([]byte(mock)))
	if err != nil {
		t.Fatal(err)
	}

	result, err := bodyToArrayWithIDs(req)
	if err != nil {
		t.Fatal(err)
	}

	expected := 4
	if len(result) != expected {
		t.Errorf("Expected number of id's to be 4 but got %v", len(result))
	}
}

// Get an user by it's index.
func TestGetUserByIndex(t *testing.T) {
	// Mock user object
	user := getTestUser()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "email"}).AddRow(user.ID, user.Username, timeNow, user.Email)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	// Router
	r := mux.NewRouter()
	r.Handle("/user/{id}", negroni.New(
		negroni.HandlerFunc(UserByIndexHandler(db)),
	))

	// Server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Do Request
	url := ts.URL + "/user/" + strconv.Itoa(user.ID)
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.StatusCode)
	}
}

// Test IPC for getting usernames based on a list of ID's
func TestGetUsernamesHandler(t *testing.T) {
	// Mock user object
	user1 := getTestUser()
	user1.ID = 1
	user1.Username = "username1"

	user2 := getTestUser()
	user2.ID = 2
	user2.Username = "username2"
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"id", "username"}).AddRow(user1.ID, user1.Username).AddRow(user2.ID, user2.Username)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs().WillReturnRows(selectByIDRows)

	type Req struct {
		Requests []*sharedModels.GetUsernamesRequest `json:"requests"`
	}
	jsonObject := &Req{}
	jsonObject.Requests = append(jsonObject.Requests, &sharedModels.GetUsernamesRequest{
		ID: user1.ID,
	})
	jsonObject.Requests = append(jsonObject.Requests, &sharedModels.GetUsernamesRequest{
		ID: user2.ID,
	})

	json, _ := json.Marshal(jsonObject)
	req, err := http.NewRequest("POST", "http://localhost/ipc/usernames", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler := GetUsernamesHandler(db)
	handler(res, req, nil)

	// Make sure database expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response expectations are met
	expected := `{"usernames":[{"id":1,"username":"username1"},{"id":2,"username":"username2"}]}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}

// In this test we'll login as user1 and we try to change user2. This is not allowed therefore we expect an error.
func TestTryUpdateOtherUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user1 := getTestUser()
	user1.ID = 1
	user1.Username = "user1"

	user2 := getTestUser()
	user2.ID = 2
	user2.Username = "user2"

	tokenString := getTokenString(cnf, user1, t)

	json, _ := json.Marshal(user2)
	req, err := http.NewRequest(http.MethodPut, "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := UpdateUserHandler(nil, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	expected := `{"message":"you can only change your own user object"}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
	if res.Result().StatusCode != 400 {
		t.Errorf("Expected statuscode to be 400 but got %v", res.Result().StatusCode)
	}
}

// In this test we'll login as user1 and we try to delete user2. This is not allowed therefore we expect an error.
func TestTryDeleteOtherUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user1 := getTestUser()
	user1.ID = 1
	user1.Username = "user1"

	user2 := getTestUser()
	user2.ID = 2
	user2.Username = "user2"

	tokenString := getTokenString(cnf, user1, t)

	json, _ := json.Marshal(user2)
	req, err := http.NewRequest(http.MethodDelete, "http://localhost/users?token="+tokenString, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := UpdateUserHandler(nil, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	expected := `{"message":"you can only change your own user object"}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
	if res.Result().StatusCode != 400 {
		t.Errorf("Expected statuscode to be 400 but got %v", res.Result().StatusCode)
	}
}

func getTestUser() *models.UserCreate {
	user := &models.UserCreate{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"
	user.CreatedAt = time.Now()
	user.Hash = "TempFakeHash"
	return user
}

func getTokenString(cnf config.Config, user *models.UserCreate, t *testing.T) string {
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	secretKey := cnf.SecretKey
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Error(err)
		return ""
	}
	return tokenString
}
