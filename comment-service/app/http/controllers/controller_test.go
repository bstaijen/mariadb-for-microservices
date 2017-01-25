package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	"github.com/buger/jsonparser"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

func TestCreateHandler(t *testing.T) {

	// MOCK SERVER
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// MOCK RESPONSE
		if r.URL.String() == "/ipc/usernames" {
			users := make([]*sharedModels.GetUsernamesResponse, 0)
			users = append(users, &sharedModels.GetUsernamesResponse{
				ID:       1,
				Username: "mockuser",
			})

			type Resp struct {
				Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
			}
			util.SendOK(w, &Resp{Usernames: users})
		}
	}))
	defer ts.Close()

	// Test Comment
	comment := &models.CommentCreate{}
	comment.Comment = "comment"
	comment.PhotoID = 5
	comment.UserID = 9

	// Prepare request. (does not go to mock server, doens't go to any server at all)
	json, err := json.Marshal(comment)
	req, err := http.NewRequest("POST", ts.URL+"/comment", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expectation: insert into database
	mock.ExpectExec("INSERT INTO comments").WithArgs(comment.UserID, comment.PhotoID, comment.Comment).WillReturnResult(sqlmock.NewResult(1, 1))

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(1, comment.UserID, comment.PhotoID, comment.Comment, timeNow)
	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs(1).WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}
	cnf.ProfileServiceBaseurl = ts.URL + "/"

	// Invoke handler
	handler := CreateHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)

		t.Errorf(res.Body.String())
	}

}

func TestListCommentsHandler(t *testing.T) {
	// Test Comment
	comment := &models.CommentCreate{}
	comment.Comment = "comment"
	comment.PhotoID = 5
	comment.UserID = 9

	// MOCK SERVER
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// MOCK RESPONSE
		if r.URL.String() == "/ipc/usernames" {
			users := make([]*sharedModels.GetUsernamesResponse, 0)
			users = append(users, &sharedModels.GetUsernamesResponse{
				ID:       1,
				Username: "mockuser",
			})

			type Resp struct {
				Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
			}
			util.SendOK(w, &Resp{Usernames: users})
		}
	}))
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL+"/comment?photoID=5", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(1, comment.UserID, comment.PhotoID, comment.Comment, timeNow)
	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs(5, 1, 10).WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}
	cnf.ProfileServiceBaseurl = ts.URL + "/"
	handler := ListCommentsHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)

		t.Errorf(res.Body.String())
	}
}

func TestGetCommentCountHandler(t *testing.T) {

	// MOCK SERVER
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// MOCK RESPONSE
		if r.URL.String() == "/ipc/usernames" {
			users := make([]*sharedModels.GetUsernamesResponse, 0)
			users = append(users, &sharedModels.GetUsernamesResponse{
				ID:       1,
				Username: "mockuser",
			})

			type Resp struct {
				Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
			}
			util.SendOK(w, &Resp{Usernames: users})
		}
	}))
	defer ts.Close()

	body := []byte(`{ "requests":[{"photo_id":1} ,{"photo_id":2},{"photo_id":3}, {"photo_id":4} ]}`)

	req, err := http.NewRequest("POST", ts.URL+"/comment?photoID=5", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "count(*)"}).AddRow(5, 10).AddRow(6, 11)
	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs().WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}
	cnf.ProfileServiceBaseurl = ts.URL + "/"
	handler := GetCommentCountHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)

		t.Errorf(res.Body.String())
	}

	expected := `{"result":[{"photo_id":5,"count":10},{"photo_id":6,"count":11}]}`
	actual := res.Body.String()
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestGetLastTenHandler(t *testing.T) {
	// MOCK SERVER
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// MOCK RESPONSE
		if r.URL.String() == "/ipc/usernames" {
			users := make([]*sharedModels.GetUsernamesResponse, 0)
			users = append(users, &sharedModels.GetUsernamesResponse{
				ID:       1,
				Username: "mockuser",
			})

			type Resp struct {
				Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
			}
			util.SendOK(w, &Resp{Usernames: users})
		}
	}))
	defer ts.Close()

	body := []byte(`{ "requests":[{"photo_id":1} ,{"photo_id":2},{"photo_id":3}, {"photo_id":4} ]}`)

	req, err := http.NewRequest("POST", ts.URL+"/comment?photoID=5", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Add Rows
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "photo_id", "comment", "createdAt"}).AddRow(1, 1, 5, "comment", timeNow)

	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs().WillReturnRows(selectByIDRows)

	// Mock config
	cnf := config.Config{}
	cnf.ProfileServiceBaseurl = ts.URL + "/"
	handler := GetLastTenHandler(db, cnf)
	handler(res, req, nil)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)

		t.Errorf(res.Body.String())
	}

	// Compare results
	objects := make([]*sharedModels.CommentResponse, 0)
	jsonparser.ArrayEach([]byte(res.Body.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		comment := &sharedModels.CommentResponse{}
		json.Unmarshal(value, comment)
		objects = append(objects, comment)
	}, "comments")

	returnComment := objects[0]
	expectedID := 1
	actualID := returnComment.ID
	if expectedID != actualID {
		t.Fatalf("Expected %v got %v", expectedID, actualID)
	}

	actualPhotoID := returnComment.PhotoID
	expectedPhotoID := 5
	if expectedPhotoID != actualPhotoID {
		t.Fatalf("Expected %v got %v", expectedPhotoID, actualPhotoID)
	}

	actualUserID := returnComment.UserID
	expectedUserID := 1
	if expectedUserID != actualUserID {
		t.Fatalf("Expected %v got %v", expectedUserID, actualUserID)
	}

	actualUsername := returnComment.Username
	expectedUsername := "mockuser"
	if expectedUsername != actualUsername {
		t.Fatalf("Expected %v got %v", expectedUsername, actualUsername)
	}

	actualComment := returnComment.Comment
	expectedComment := "comment"
	if expectedComment != actualComment {
		t.Fatalf("Expected %v got %v", expectedComment, actualComment)
	}

}

// This test is only testing of the method getUsernames can be called and if it correctly parses the data back. It is not meant to test whether the whole IPC-call is performed correctly.
func TestGetUsernames(t *testing.T) {

	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/ipc/usernames" {
			users := make([]*sharedModels.GetUsernamesResponse, 0)
			users = append(users, &sharedModels.GetUsernamesResponse{
				ID:       19,
				Username: "mockuser19",
			})
			users = append(users, &sharedModels.GetUsernamesResponse{
				ID:       54,
				Username: "mockuser54",
			})

			type Resp struct {
				Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
			}
			util.SendOK(w, &Resp{Usernames: users})
		}
	}))
	defer ts.Close()

	// Mock config
	cnf := config.Config{}
	cnf.ProfileServiceBaseurl = ts.URL + "/"

	// Mocking list with desired ID's.
	identifiers := make([]*sharedModels.GetUsernamesRequest, 0)
	identifiers = append(identifiers, &sharedModels.GetUsernamesRequest{
		ID: 19,
	})
	identifiers = append(identifiers, &sharedModels.GetUsernamesRequest{
		ID: 54,
	})

	result := getUsernames(cnf, identifiers)
	if len(result) == 2 {
		// OK
	} else {
		t.Errorf("Expected 2 rows to return but got %v", len(result))
	}
}
