package db

import (
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetUserByID(t *testing.T) {
	user := getTestUserForCreation()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expectation: select user where id == user.ID
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(user.ID, user.Username, timeNow, user.Hash, user.Email)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.ID).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := GetUserByID(db, user.ID); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := getTestUserForCreation()

	// Expected rows
	rows := sqlmock.NewRows([]string{"count(*)"})

	// define expectations
	// Expectation: check for unique username
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Username).WillReturnRows(rows)

	// Expectation: check for unique email
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Email).WillReturnRows(rows)

	// Expectation: insert into database
	mock.ExpectExec("INSERT INTO users").WithArgs(user.Username, user.Email, user.Hash).WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the method
	if _, err := CreateUser(db, user); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateUser
func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := getTestUser()

	// Expectation: insert into database
	mock.ExpectExec("UPDATE users SET").WithArgs(user.Username, user.Email, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the method
	if _, err := UpdateUser(db, user); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDeleteUser
func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := getTestUser()

	// Expectation: insert into database
	mock.ExpectExec("DELETE from users WHERE").WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the method
	if _, err := DeleteUser(db, user); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetUsers
func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user1 := getTestUser()
	user2 := getTestUser()
	user2.ID = 2

	// Expectation: insert into database
	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "email", "createdAt"}).AddRow(user1.ID, user1.Username, user1.Email, timeNow).AddRow(user2.ID, user2.Username, user2.Email, timeNow)
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := GetUsers(db); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetUsernames
func TestGetUsernames(t *testing.T) {
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

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"id", "username"}).AddRow(user1.ID, user1.Username).AddRow(user2.ID, user2.Username)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs().WillReturnRows(selectByIDRows)

	ids := []*sharedModels.GetUsernamesRequest{}
	ids = append(ids, &sharedModels.GetUsernamesRequest{
		ID: user1.ID,
	})
	ids = append(ids, &sharedModels.GetUsernamesRequest{
		ID: user2.ID,
	})

	// Execute the method
	if _, err := GetUsernames(db, ids); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestQueryBuilder
func TestQueryBuilder(t *testing.T) {
	user1 := getTestUser()
	user1.ID = 1
	user1.Username = "username1"

	user2 := getTestUser()
	user2.ID = 2
	user2.Username = "username2"

	ids := []*sharedModels.GetUsernamesRequest{}

	// If there's no input, then there's no query.
	actual := inQueryBuilder(ids)
	expected := ""
	if expected != actual {
		t.Errorf("we are expecting '%s' instead we got '%s'", actual, expected)
	}

	ids = append(ids, &sharedModels.GetUsernamesRequest{
		ID: user1.ID,
	})
	ids = append(ids, &sharedModels.GetUsernamesRequest{
		ID: user2.ID,
	})

	// If there's input. Then we'd expect a query.
	actual = inQueryBuilder(ids)
	expected = "SELECT id, username FROM users WHERE id IN(1,2)"
	if expected != actual {
		t.Errorf("we are expecting '%s' instead we got '%s'", actual, expected)
	}
}

// TestUniqueEmail
func TestUniqueEmail(t *testing.T) {
	// Mock user object
	user1 := getTestUserForCreation()
	user1.ID = 1
	user1.Username = "username1"

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	rows := sqlmock.NewRows([]string{"count(*)"})
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user1.Username).WillReturnRows(rows)

	// We add 1 row so the code will think the database returned a row (which means there's a duplicate).
	rowsWithData := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user1.Email).WillReturnRows(rowsWithData)

	// Execute the method
	if _, err := CreateUser(db, user1); err != nil {
		expected := ErrEmailIsNotUnique.Error()
		actual := err.Error()
		if expected != actual {
			t.Errorf("we are expecting error '%s' instead we got '%s'", ErrEmailIsNotUnique, err)
		}
	} else {
		t.Errorf("we are expecting error '%s' instead the code threw no error.", ErrEmailIsNotUnique)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUniqueUsername
func TestUniqueUsername(t *testing.T) {
	// Mock user object
	user1 := getTestUserForCreation()
	user1.ID = 1
	user1.Username = "username1"

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	// We add 1 row so the code will think the database returned a row (which means there's a duplicate).
	rowsWithData := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user1.Username).WillReturnRows(rowsWithData)

	// Execute the method
	if _, err := CreateUser(db, user1); err != nil {
		expected := ErrUsernameIsNotUnique.Error()
		actual := err.Error()
		if expected != actual {
			t.Errorf("we are expecting error '%s' instead we got '%s'", ErrUsernameIsNotUnique, err)
		}
	} else {
		t.Errorf("we are expecting error '%s' instead the code threw no error.", ErrUsernameIsNotUnique)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getTestUserForCreation() *models.UserCreate {
	user := &models.UserCreate{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"
	user.CreatedAt = time.Now()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Hash = string(hash)
	return user
}

func getTestUser() *models.UserResponse {
	user := &models.UserResponse{}
	user.ID = 1
	user.Email = "username@example.com"
	user.Username = "username"
	user.CreatedAt = time.Now()
	return user
}
