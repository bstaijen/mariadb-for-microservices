package db

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define user
	user := &models.User{}
	user.Email = "username@example.com"
	user.Password = "password"
	user.Username = "username"
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Hash = string(hash)

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
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

// TestUniqueEmail
// TestUniqueUsername
// TestUpdate
// TestGetAll
// TestGetByID
// TestDeleteUsers
// TestQueryBuilder
