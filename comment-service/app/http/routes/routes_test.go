package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// test post comment
func TestPostComment(t *testing.T) {

	// Test Comment
	comment := &models.CommentCreate{}
	comment.Comment = "comment"
	comment.PhotoID = 5
	comment.UserID = 9

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
		} else {
			util.SendBadRequest(w, errors.New("Not implemented"))
		}
	}))
	defer ts.Close()

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

	cnf := config.Config{}
	cnf.ProfileServiceBaseurl = ts.URL + "/"

	json, _ := json.Marshal(comment)
	res := doRequest(db, cnf, "POST", ts.URL+"/comments", bytes.NewBuffer(json), t)
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

// test get comments on photo
func TestGetCommentsFromPhoto(t *testing.T) {
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
				ID:       5,
				Username: "mockuser",
			})

			type Resp struct {
				Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
			}
			util.SendOK(w, &Resp{Usernames: users})
		} else {
			util.SendBadRequest(w, errors.New("Not implemented"))
		}
	}))
	defer ts.Close()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(1, comment.UserID, comment.PhotoID, comment.Comment, timeNow)
	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs(5, 1, 10).WillReturnRows(selectByIDRows)

	cnf := config.Config{}
	cnf.ProfileServiceBaseurl = ts.URL + "/"

	res := doRequest(db, cnf, "GET", ts.URL+"/comments?photoID=5", bytes.NewBuffer([]byte("")), t)

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

// test get list 10 -> ipc
func TestIPCGetLast10(t *testing.T) {
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
		} else {
			util.SendBadRequest(w, errors.New("Not implemented"))
		}
	}))
	defer ts.Close()

	body := []byte(`{ "requests":[{"photo_id":1} ,{"photo_id":2},{"photo_id":3}, {"photo_id":4} ]}`)

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

	res := doRequest(db, cnf, "GET", ts.URL+"/ipc/getLast10", bytes.NewBuffer(body), t)

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v", res.Result().StatusCode)
		t.Errorf(res.Body.String())
	}

	// Compare results
	type Collection struct {
		Objects []*sharedModels.CommentResponse `json:"comments"`
	}
	col := &Collection{}
	col.Objects = make([]*sharedModels.CommentResponse, 0)
	err = util.ResponseRecorderJSONToObject(res, &col)
	if err != nil {
		log.Fatal(err)
	}

	returnComment := col.Objects[0]
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

// test get comment count -> ipc
func TestIPCGetCommentCount(t *testing.T) {

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
		} else {
			util.SendBadRequest(w, errors.New("Not implemented"))
		}
	}))
	defer ts.Close()

	body := []byte(`{ "requests":[{"photo_id":1} ,{"photo_id":2},{"photo_id":3}, {"photo_id":4} ]}`)

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

	res := doRequest(db, cnf, "GET", ts.URL+"/ipc/getCount", bytes.NewBuffer(body), t)

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
