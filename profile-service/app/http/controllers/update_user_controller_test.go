package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestUpdateUser(t *testing.T) {
	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"

	user := &models.User{}
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

	handler := controllers.UpdateUserHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	}
}
