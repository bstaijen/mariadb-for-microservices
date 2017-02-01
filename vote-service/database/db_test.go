package db

import (
	"testing"

	"encoding/json"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreate(t *testing.T) {
	toCreate := &sharedModels.VoteCreateRequest{}
	toCreate.Downvote = false
	toCreate.Upvote = true
	toCreate.PhotoID = 5
	toCreate.UserID = 9

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM votes WHERE").WithArgs(toCreate.UserID, toCreate.PhotoID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO votes").WithArgs(toCreate.UserID, toCreate.PhotoID, toCreate.Upvote, toCreate.Downvote).WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the method
	if err := Create(db, toCreate); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestVoteCount(t *testing.T) {

	list := make([]*sharedModels.VoteCountRequest, 0)
	list = append(list, &sharedModels.VoteCountRequest{
		PhotoID: 1,
	})
	list = append(list, &sharedModels.VoteCountRequest{
		PhotoID: 2,
	})

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "sum(upvote)", "sum(downvote)"}).AddRow(1, 6, 5).AddRow(2, 8, 4)
	mock.ExpectQuery("SELECT (.+) FROM votes WHERE").WithArgs().WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := VoteCount(db, list); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHasVoted(t *testing.T) {
	list := make([]*sharedModels.HasVotedRequest, 0)
	list = append(list, &sharedModels.HasVotedRequest{
		PhotoID: 1,
		UserID:  9,
	})
	list = append(list, &sharedModels.HasVotedRequest{
		PhotoID: 2,
		UserID:  11,
	})

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"user_id", "photo_id", "upvote", "downvote"}).AddRow(9, 1, true, false).AddRow(11, 2, false, true)
	mock.ExpectQuery("SELECT (.+) FROM votes WHERE").WithArgs().WillReturnRows(selectByIDRows)

	// Execute the method
	if _, err := HasVoted(db, list); err != nil {
		t.Errorf("there was an unexpected error: %s", err)
	}

	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTopRatedTimeline(t *testing.T) {

	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "total_upvote", "total_downvote", "difference"}).AddRow(1, 4, 5, 1).AddRow(2, 8, 3, 5)
	mock.ExpectQuery("SELECT (.+) FROM votes").WithArgs(1, 10).WillReturnRows(selectByIDRows)

	// Execute the method
	if result, err := GetTopRatedTimeline(db, 1, 10); err != nil {
		t.Errorf("there was an unexpected error: %s", err)

	} else {
		json, err := json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(json))
	}
	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestGetHotTimeline(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	// Database expectations
	selectByIDRows := sqlmock.NewRows([]string{"photo_id", "total_upvote", "total_downvote", "difference"}).AddRow(1, 4, 5, 1).AddRow(2, 8, 3, 5)
	mock.ExpectQuery("SELECT (.+) FROM votes").WithArgs(1, 10).WillReturnRows(selectByIDRows)

	// Execute the method
	if result, err := GetHotTimeline(db, 1, 10); err != nil {
		t.Errorf("there was an unexpected error: %s", err)

	} else {
		json, err := json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(json))
	}
	// Make sure expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
