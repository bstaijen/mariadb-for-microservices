package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/config"
	jwt "github.com/dgrijalva/jwt-go"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

func TestOPTIONSVotes(t *testing.T) {
	// Router
	r := InitRoutes(nil, config.Config{})
	res := httptest.NewRecorder()

	// Do Request
	req, err := http.NewRequest(http.MethodOptions, "/votes", nil)
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

func TestPOSTVotes(t *testing.T) {
	vote := getTestVote()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM votes WHERE").WithArgs(vote.UserID, vote.PhotoID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO votes").WithArgs(vote.UserID, vote.PhotoID, vote.Upvote, vote.Downvote).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock config
	cnf := config.Config{}

	// Get token string
	tokenString := getTokenString(cnf, vote.UserID, t)

	json, _ := json.Marshal(vote)
	res := doRequest(db, cnf, http.MethodPost, "/votes?token="+tokenString, bytes.NewBuffer(json), t)

	t.Log(res.Body.String())

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

func TestIPCGetTopRated(t *testing.T) {

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "total_upvote", "total_downvote", "difference"}).AddRow(1, 4, 5, 1).AddRow(2, 8, 3, 5)
	mock.ExpectQuery("SELECT (.+) FROM votes").WithArgs(0, 10).WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}

	res := doRequest(db, cnf, http.MethodGet, "/ipc/toprated", bytes.NewBuffer([]byte("")), t)

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

func TestIPCGetHot(t *testing.T) {

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "total_upvote", "total_downvote", "difference"}).AddRow(1, 4, 5, 1).AddRow(2, 8, 3, 5)
	mock.ExpectQuery("SELECT (.+) FROM votes").WithArgs(0, 10).WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}

	res := doRequest(db, cnf, http.MethodGet, "/ipc/hot", bytes.NewBuffer([]byte("")), t)

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

func TestIPCGetCount(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "sum(upvote)", "sum(downvote)"}).AddRow(1, 6, 5).AddRow(2, 8, 4)
	mock.ExpectQuery("SELECT (.+) FROM votes WHERE").WithArgs().WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}

	body := []byte(`{ "requests":[{"photo_id":1} ,{"photo_id":2},{"photo_id":3}, {"photo_id":4} ]}`)

	res := doRequest(db, cnf, http.MethodGet, "/ipc/count", bytes.NewBuffer(body), t)

	t.Log(res.Body.String())

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

func TestIPCGetVoted(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"user_id", "photo_id", "sum(upvote)", "sum(downvote)"}).AddRow(1, 1, false, true).AddRow(2, 2, true, false)
	mock.ExpectQuery("SELECT (.+) FROM votes WHERE").WithArgs().WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}

	body := []byte(`{ "requests":[{"photo_id":1} ,{"photo_id":2},{"photo_id":3}, {"photo_id":4} ]}`)

	res := doRequest(db, cnf, http.MethodGet, "/ipc/voted", bytes.NewBuffer(body), t)

	t.Log(res.Body.String())

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

func getTestVote() *sharedModels.VoteCreateRequest {
	vote := &sharedModels.VoteCreateRequest{}
	vote.Downvote = false
	vote.Upvote = true
	vote.UserID = 5
	vote.PhotoID = 9
	return vote
}

func getTokenString(cnf config.Config, userID int, t *testing.T) string {
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
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
