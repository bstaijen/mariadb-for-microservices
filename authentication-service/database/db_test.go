package db

import (
	"testing"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/models"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := getTestUser()

	// Expected rows
	timeNow := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(user.ID, user.Username, timeNow, "TestUserHash", user.Email)

	// define expectations
	// Expectation: check for unique username
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(user.Username).WillReturnRows(rows)

	// Execute the method
	if _, err := GetUserByUsername(db, user.Username); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
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
