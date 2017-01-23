package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetUserByIndex(t *testing.T) {
	// Mock user object
	user := &models.User{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(user.ID, user.Username, timeNow, "PasswordHashPlaceHolder", user.Email)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	// Router
	r := mux.NewRouter()
	r.Handle("/user/{id}", negroni.New(
		negroni.HandlerFunc(controllers.UserByIndexHandler(db)),
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

/*func TestGetUsernamesHandler(t *testing.T) {
	//

	// Mock user object
	user1 := &models.User{}
	user1.ID = 1
	user1.Username = "username1"

	user2 := &models.User{}
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

	handler := controllers.GetUsernamesHandler(db)
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
*/
