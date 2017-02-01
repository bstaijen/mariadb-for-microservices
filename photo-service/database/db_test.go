package db

import (
	"testing"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestInsertPhoto(t *testing.T) {

	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expectation: insert into database
	mock.ExpectExec("INSERT INTO photos").WithArgs(photo.UserID, photo.Filename, photo.Title, photo.ContentType, photo.Image).WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the method
	if err := InsertPhoto(db, photo); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListImagesByUserID(t *testing.T) {
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos WHERE").WithArgs(1).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := ListImagesByUserID(db, 1); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListIncoming(t *testing.T) {
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos").WithArgs(1, 10).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := ListIncoming(db, 1, 10); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPhotoByFilename(t *testing.T) {
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos").WithArgs(photo.Filename).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := GetPhotoByFilename(db, photo.Filename); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPhotoById(t *testing.T) {
	photo := &models.CreatePhoto{}
	photo.ContentType = "image/png"
	photo.Filename = "test.png"
	photo.Image = []byte(`ABCDEFGHIJ`)
	photo.Title = "Test image"
	photo.UserID = 1

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeNow := time.Now().UTC()
	selectByIDRows := sqlmock.NewRows([]string{"id", "user_id", "filename", "title", "createdAt", "contentType", "photo"}).AddRow(1, photo.UserID, photo.Filename, photo.Title, timeNow, photo.ContentType, photo.Image)
	mock.ExpectQuery("SELECT (.+) FROM photos").WithArgs(1).WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := GetPhotoById(db, 1); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
