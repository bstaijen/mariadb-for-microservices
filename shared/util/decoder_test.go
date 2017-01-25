package util

import (
	"bytes"
	"net/http"
	"testing"
)

func TestRequestJSONToObject(t *testing.T) {

	mock := []byte(`{"username":"test"}`)
	req, err := http.NewRequest("GET", "http://localhost/", bytes.NewBuffer([]byte(mock)))
	if err != nil {
		t.Fatal(err)
	}

	type User struct {
		Username string `json:"username"`
	}
	user := &User{}
	err = RequestToJSON(req, user)
	if err != nil {
		t.Fatal(err)
	}

	expected := "test"
	actual := user.Username
	if expected != actual {
		t.Fatalf("expected %s instead got %s", expected, actual)
	}
}

func TestRequestJSONToArray(t *testing.T) {
	mock := []byte(`[{"username":"test"},{"username":"test2"}]`)
	req, err := http.NewRequest("GET", "http://localhost/", bytes.NewBuffer([]byte(mock)))
	if err != nil {
		t.Fatal(err)
	}

	type User struct {
		Username string `json:"username"`
	}

	user := make([]User, 0)
	err = RequestToJSON(req, &user)
	if err != nil {
		t.Fatal(err)
	}

	expected := "test"
	actual := user[0].Username
	if expected != actual {
		t.Fatalf("expected %s instead got %s", expected, actual)
	}

	expected = "test2"
	actual = user[1].Username
	if expected != actual {
		t.Fatalf("expected %s instead got %s", expected, actual)
	}
}

func TestRequestBadJSON(t *testing.T) {
	mock := []byte(`{`)
	req, err := http.NewRequest("GET", "http://localhost/", bytes.NewBuffer([]byte(mock)))
	if err != nil {
		t.Fatal(err)
	}

	type User struct {
		Username string `json:"username"`
	}
	user := &User{}
	err = RequestToJSON(req, user)
	if err == nil {
		t.Fatal("We expected to get an error. Instead we got nothing.")
	}

	expected := "unexpected EOF"
	actual := err.Error()
	if expected != actual {
		t.Fatalf("expected %s instead got %s", expected, actual)
	}
}

func TestRequestNoJSON(t *testing.T) {
	mock := []byte(``)
	req, err := http.NewRequest("GET", "http://localhost/", bytes.NewBuffer([]byte(mock)))
	if err != nil {
		t.Fatal(err)
	}

	type User struct {
		Username string `json:"username"`
	}
	user := &User{}
	err = RequestToJSON(req, user)
	if err == nil {
		t.Fatal("We expected to get an error. Instead we got nothing.")
	}

	expected := "EOF"
	actual := err.Error()
	if expected != actual {
		t.Fatalf("expected %s instead got %s", expected, actual)
	}
}
