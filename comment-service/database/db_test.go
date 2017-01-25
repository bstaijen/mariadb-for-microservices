package db

import (
	"testing"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/models"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

func TestCreate(t *testing.T) {
	// Test Comment
	comment := &models.CommentCreate{}
	comment.Comment = "comment"
	comment.PhotoID = 5
	comment.UserID = 9

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

	// Execute the method
	if _, err := Create(db, comment); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetCommentByID(t *testing.T) {
	// Test Comment
	comment := &models.CommentCreate{}
	comment.Comment = "comment"
	comment.PhotoID = 5
	comment.UserID = 9

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(1, comment.UserID, comment.PhotoID, comment.Comment, timeNow)
	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs(comment.PhotoID).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := GetCommentByID(db, comment.PhotoID); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetComments(t *testing.T) {
	// Test Comment
	comment := &models.CommentCreate{}
	comment.Comment = "comment"
	comment.PhotoID = 5
	comment.UserID = 9

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "username", "createdAt", "password", "email"}).AddRow(1, comment.UserID, comment.PhotoID, comment.Comment, timeNow)
	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs(5, 1, 10).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := GetComments(db, comment.PhotoID, 1, 10); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestCommentCount(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "count(*)"}).AddRow(5, 10).AddRow(6, 11)
	mock.ExpectQuery("SELECT (.+) FROM comments WHERE").WithArgs().WillReturnRows(selectByIDRows)

	objects := make([]*sharedModels.CommentCountRequest, 0)
	objects = append(objects, &sharedModels.CommentCountRequest{
		PhotoID: 5,
	})

	// Execute the method
	if _, err := GetCommentCount(db, objects); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetLastTenComments(t *testing.T) {
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

	objects := make([]*sharedModels.CommentRequest, 0)
	objects = append(objects, &sharedModels.CommentRequest{
		PhotoID: 5,
	})

	// Execute the method
	if _, err := GetLastTenComments(db, objects); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
