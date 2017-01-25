package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/config"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestPOSTTokenAuth(t *testing.T) {
	user := getTestUser()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(user.ID, user.Username, timeNow, hash, user.Email)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Username).WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "ABC"

	json, _ := json.Marshal(user)
	res := doRequest(db, cnf, http.MethodPost, "/token-auth", bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
		t.Errorf(res.Body.String())
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func doRequest(db *sql.DB, cnf config.Config, method string, url string, body *bytes.Buffer, t *testing.T) *httptest.ResponseRecorder {
	r := InitRoutes(db, cnf)
	res := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(res, req)
	return res
}

func getTestUser() *models.User {
	user := &models.User{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"
	user.CreatedAt = time.Now()
	return user
}

func getTokenString(cnf config.Config, user *models.User, t *testing.T) string {
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
