package routes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"strconv"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	jwt "github.com/dgrijalva/jwt-go"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestHash struct{}

func (a TestHash) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func TestOPTIONSUsers(t *testing.T) {
	// Router
	r := InitRoutes(nil, config.Config{})
	res := httptest.NewRecorder()

	// Do Request
	req, err := http.NewRequest(http.MethodOptions, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(res, req)

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func TestPUTUsers(t *testing.T) {

	user := getTestUser()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE users SET").WithArgs(user.Username, user.Email, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "ABC"

	// Get token string
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	res := doRequest(db, cnf, http.MethodPut, "/users?token="+tokenString, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func TestDELETEUsers(t *testing.T) {
	user := getTestUser()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE from users").WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "ABC"

	// Get token string
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	res := doRequest(db, cnf, http.MethodDelete, "/users?token="+tokenString, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func TestPOSTUsers(t *testing.T) {
	user := getTestUser()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expected rows
	rows := sqlmock.NewRows([]string{"count (*)"})

	// Expectation: check for unique username
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Username).WillReturnRows(rows)

	// Expectation: check for unique email
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Email).WillReturnRows(rows)

	// Expectation: insert into database
	mock.ExpectExec("INSERT INTO users").WithArgs(user.Username, user.Email, TestHash{}).WillReturnResult(sqlmock.NewResult(1, 1))

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "email"}).AddRow(user.ID, user.Username, timeNow, user.Email)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "ABC"

	// Get token string
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	res := doRequest(db, cnf, http.MethodPost, "/users?token="+tokenString, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Logf(string(res.Body.Bytes()))
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func TestGETUserByID(t *testing.T) {
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

	// Mock config
	cnf := config.Config{}
	cnf.SecretKey = "ABC"

	// Get token string
	tokenString := getTokenString(cnf, user, t)

	json, _ := json.Marshal(user)
	url := "/user/" + strconv.Itoa(user.ID) + "?token=" + tokenString
	res := doRequest(db, cnf, http.MethodGet, url, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
	} else {
		t.Logf("Result statuscode %v. (As expected)", res.Result().StatusCode)
	}
}

func TestIPCGetUsernames(t *testing.T) {
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

	cnf := config.Config{}

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
	url := "/ipc/usernames"
	res := doRequest(db, cnf, http.MethodGet, url, bytes.NewBuffer(json), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Make sure response statuscode expectation is met
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
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
