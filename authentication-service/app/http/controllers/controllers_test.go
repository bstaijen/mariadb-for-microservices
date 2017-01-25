package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/config"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestLoginHandlerWithoutUsername(t *testing.T) {
	user := &models.User{}

	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	// Mock database
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "ABC"

	handler := LoginHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Please provide username and password in the body\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}

}

func TestLoginHandlerWithoutPassword(t *testing.T) {
	user := &models.User{}
	user.Username = "username"
	json, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "http://localhost/users", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	// Mock database
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "ABC"

	handler := LoginHandler(db, cnf)
	handler(res, req, nil)

	actual := res.Body.String()
	expected := "{\"message\":\"Please provide username and password in the body\"}"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}
