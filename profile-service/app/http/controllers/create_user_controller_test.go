package controllers_test

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestHash struct{}

func (a TestHash) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func TestCreateUser(t *testing.T) {

	user := &models.User{}
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
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(user.ID, user.Username, timeNow, "PasswordHashPlaceHolder", user.Email)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	handler := controllers.CreateUserHandler(db)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	// Make sure response is alright
	responseUser := &models.User{}
	err = toJSON(res.Body, responseUser)
	if err != nil {
		t.Fatal(errors.New("Bad json"))
		return
	}
	if responseUser.ID < 1 {
		t.Errorf("Expected user ID greater than 0 but got %v", responseUser.ID)
	}

	if responseUser.Username != user.Username {
		t.Errorf("Expected username to be %v but got %v", user.Username, responseUser.Username)
	}
	if responseUser.Email != user.Email {
		t.Errorf("Expected username to be %v but got %v", user.Email, responseUser.Email)
	}

	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}

func TestBadJson(t *testing.T) {
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
	handler := controllers.CreateUserHandler(db)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Bad json\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestCreateUserWithoutUsername(t *testing.T) {
	user := &models.User{}
	//user.Username = "username"

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

	handler := controllers.CreateUserHandler(db)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Username is too short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestCreateUserWithoutPassword(t *testing.T) {
	user := &models.User{}
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

	handler := controllers.CreateUserHandler(db)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Password is to short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestCreateUserWithoutEmail(t *testing.T) {
	user := &models.User{}
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

	handler := controllers.CreateUserHandler(db)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Email address is to short\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func toJSON(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		fmt.Printf("json decoder error occured: %v \n", err.Error())
		return err
	}
	return nil
}
