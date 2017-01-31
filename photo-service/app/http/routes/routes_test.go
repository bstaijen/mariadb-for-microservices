package routes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	jwt "github.com/dgrijalva/jwt-go"
)

type TestFilename struct{}

func (a TestFilename) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func TestOPTIONSImage(t *testing.T) {
	// Router
	r := InitRoutes(nil, config.Config{})
	res := httptest.NewRecorder()

	// Do Request
	req, err := http.NewRequest(http.MethodOptions, "/image", nil)
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

func TestPostImage(t *testing.T) {

	// Expected photo
	photo := &models.CreatePhoto{}
	photo.ContentType = "application/octet-stream" // this is not real
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "TestTitle"
	photo.UserID = 1

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.SendBadRequest(w, errors.New("Not implemented"))
	}))
	defer ts.Close()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expectation: insert into database
	mock.ExpectExec("INSERT INTO photos").WithArgs(photo.UserID, TestFilename{}, photo.Title, photo.ContentType, photo.Image).WillReturnResult(sqlmock.NewResult(1, 1))

	cnf := config.Config{}
	cnf.SecretKey = "ABCDEF"
	token := getTokenString(cnf, photo.UserID, t)
	res := doPostRequest(db, cnf, ts.URL+"/image/1?title=TestTitle&token="+token, bytes.NewBuffer(photo.Image), t)

	t.Log(res.Body.String())

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

func TestListImagesFromUser(t *testing.T) {
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.SendBadRequest(w, errors.New("Not implemented"))
	}))
	defer ts.Close()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos WHERE").WithArgs(1).WillReturnRows(selectByIDRows)

	cnf := config.Config{}

	res := doRequest(db, cnf, http.MethodGet, ts.URL+"/image/1/list", bytes.NewBuffer([]byte(``)), t)

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

func TestGetTimeline(t *testing.T) {
	// /list
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.SendBadRequest(w, errors.New("Not implemented"))
	}))
	defer ts.Close()

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos").WithArgs(0, 10).WillReturnRows(selectByIDRows)

	cnf := config.Config{}

	res := doRequest(db, cnf, http.MethodGet, ts.URL+"/image/list", bytes.NewBuffer([]byte(``)), t)

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

func TestGetTopratedTimeline(t *testing.T) {

	// /toprated
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	// Mock server with fake data. We need this for our IPC to succeed
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t.Log(r.URL.String())
		if r.URL.String() == "/ipc/toprated?offset=0&rows=10" {
			photos := make([]*sharedModels.TopRatedPhotoResponse, 0)
			photos = append(photos, &sharedModels.TopRatedPhotoResponse{
				PhotoID: 1,
			})
			type Resp struct {
				Results []*sharedModels.TopRatedPhotoResponse `json:"results"`
			}
			util.SendOK(w, &Resp{Results: photos})
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
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos WHERE").WithArgs(1).WillReturnRows(selectByIDRows)

	cnf := config.Config{}
	cnf.CommentServiceBaseurl = ts.URL + "/"
	cnf.ProfileServiceBaseurl = ts.URL + "/"
	cnf.VoteServiceBaseurl = ts.URL + "/"

	res := doRequest(db, cnf, http.MethodGet, ts.URL+"/image/toprated", bytes.NewBuffer([]byte(``)), t)

	t.Log(res.Body.String())

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

func TestGetHotTimeline(t *testing.T) {
	// /toprated
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	// Mock server with fake data. We need this for our IPC to succeed
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t.Log(r.URL.String())
		if r.URL.String() == "/ipc/hot?offset=0&rows=10" {
			photos := make([]*sharedModels.TopRatedPhotoResponse, 0)
			photos = append(photos, &sharedModels.TopRatedPhotoResponse{
				PhotoID: 1,
			})
			type Resp struct {
				Results []*sharedModels.TopRatedPhotoResponse `json:"results"`
			}
			util.SendOK(w, &Resp{Results: photos})
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
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos WHERE").WithArgs(1).WillReturnRows(selectByIDRows)

	cnf := config.Config{}
	cnf.CommentServiceBaseurl = ts.URL + "/"
	cnf.ProfileServiceBaseurl = ts.URL + "/"
	cnf.VoteServiceBaseurl = ts.URL + "/"

	res := doRequest(db, cnf, http.MethodGet, ts.URL+"/image/hot", bytes.NewBuffer([]byte(``)), t)

	t.Log(res.Body.String())

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

func doPostRequest(db *sql.DB, cnf config.Config, url string, body *bytes.Buffer, t *testing.T) *httptest.ResponseRecorder {
	r := InitRoutes(db, cnf)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	res := httptest.NewRecorder()

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", "filename.png")
	if err != nil {
		fmt.Println("error writing to buffer")
		t.Fatal(err)
	}

	//iocopy fh to fileWriter
	_, err = io.Copy(fileWriter, body)
	if err != nil {
		t.Fatal(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest(http.MethodPost, url, bodyBuf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	// return response
	r.ServeHTTP(res, req)
	return res
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
