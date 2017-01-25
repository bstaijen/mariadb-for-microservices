package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {

	isHandlerFuncCalled := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHandlerFuncCalled = true
	}))
	defer ts.Close()

	isCallbackCalled := false

	err := Request("GET", ts.URL, nil, func(res *http.Response) {
		expected := true
		if isHandlerFuncCalled != expected {
			t.Errorf("Expected %v but got %v", expected, isHandlerFuncCalled)
		}

		expectedStatus := 200
		if res.StatusCode != expectedStatus {
			t.Errorf("Expected %v but got %v", expectedStatus, res.StatusCode)
		}

		isCallbackCalled = true
	})

	if err != nil {
		t.Errorf("Expected no error, instead got %v", err.Error())
	}

	expectedCb := true
	if isCallbackCalled != expectedCb {
		t.Errorf("Expected %v but got %v", expectedCb, isCallbackCalled)
	}
}

func TestBadRequest(t *testing.T) {
	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	// Do the request, but dont expect callback. Want to test proper error handling here.
	isCallbackCalled := false
	err := Request("G{}T", ts.URL, nil, func(res *http.Response) {
		isCallbackCalled = true
	})

	// We expect an error. So if we don't get one -> fail test.
	if err == nil {
		t.Error("Expected 'net/http: invalid method \"G{}T\"' error, instead got nothing")
	}

	// To make sure that the callback is not being called if there's an error.
	expectedCb := false
	if isCallbackCalled != expectedCb {
		t.Errorf("Expected %v but got %v", expectedCb, isCallbackCalled)
	}
}

func TestBadURLRequest(t *testing.T) {
	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	// Do the request, but dont expect callback. Want to test proper error handling here.
	isCallbackCalled := false
	err := Request("GET", "", nil, func(res *http.Response) {
		isCallbackCalled = true
	})

	// We expect an error. So if we don't get one -> fail test.
	if err == nil {
		t.Error("Expected 'Get : unsupported protocol scheme' error, instead got nothing")
	}

	// To make sure that the callback is not being called if there's an error.
	expectedCb := false
	if isCallbackCalled != expectedCb {
		t.Errorf("Expected %v but got %v", expectedCb, isCallbackCalled)
	}
}
