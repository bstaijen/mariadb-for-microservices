package util

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendOK(t *testing.T) {
	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendOKMessage(w, "Okay")
	}))
	defer ts.Close()

	isCallbackCalled := false
	err := Request("GET", ts.URL, nil, func(res *http.Response) {
		// Make sure code reaches the callback
		isCallbackCalled = true

		// Make sure status code is as expected
		expectedStatus := 200
		if res.StatusCode != expectedStatus {
			t.Errorf("Expected %v but got %v", expectedStatus, res.StatusCode)
		}

		// Read
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected no error, instead got %v", err.Error())
		}
		defer res.Body.Close()

		// Make sure response is as expected
		expected := "{\"message\":\"Okay\"}"
		actual := string(data)
		if expected != actual {
			t.Errorf("Expected %v but got %v", expected, actual)
		}
	})

	// We expect no error.
	if err != nil {
		t.Errorf("Expected no error, instead got %v", err.Error())
	}

	// To make sure that the callback is being called.
	expectedCb := true
	if isCallbackCalled != expectedCb {
		t.Errorf("Expected %v but got %v", expectedCb, isCallbackCalled)
	}
}

func TestSendBadRequest(t *testing.T) {
	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendBadRequest(w, errors.New("Nope"))
	}))
	defer ts.Close()

	isCallbackCalled := false
	err := Request("GET", ts.URL, nil, func(res *http.Response) {
		// Make sure code reaches the callback
		isCallbackCalled = true

		// Make sure status code is as expected
		expectedStatus := 400
		if res.StatusCode != expectedStatus {
			t.Errorf("Expected %v but got %v", expectedStatus, res.StatusCode)
		}

		// Read
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected no error, instead got %v", err.Error())
		}
		defer res.Body.Close()

		//
		expected := "{\"message\":\"Nope\"}"
		actual := string(data)
		if expected != actual {
			t.Errorf("Expected %v but got %v", expected, actual)
		}
	})

	// We expect no error.
	if err != nil {
		t.Errorf("Expected no error, instead got %v", err.Error())
	}

	// To make sure that the callback is being called.
	expectedCb := true
	if isCallbackCalled != expectedCb {
		t.Errorf("Expected %v but got %v", expectedCb, isCallbackCalled)
	}
}

func TestSendErrorMessage(t *testing.T) {
	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendErrorMessage(w, "A Error")
	}))
	defer ts.Close()

	isCallbackCalled := false
	err := Request("GET", ts.URL, nil, func(res *http.Response) {
		// Make sure code reaches the callback
		isCallbackCalled = true

		// Make sure status code is as expected
		expectedStatus := 400
		if res.StatusCode != expectedStatus {
			t.Errorf("Expected %v but got %v", expectedStatus, res.StatusCode)
		}

		// Read
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected no error, instead got %v", err.Error())
		}
		defer res.Body.Close()

		// Make sure response is as expected
		expected := "{\"message\":\"A Error\"}"
		actual := string(data)
		if expected != actual {
			t.Errorf("Expected %v but got %v", expected, actual)
		}
	})

	// We expect no error.
	if err != nil {
		t.Errorf("Expected no error, instead got %v", err.Error())
	}

	// To make sure that the callback is being called.
	expectedCb := true
	if isCallbackCalled != expectedCb {
		t.Errorf("Expected %v but got %v", expectedCb, isCallbackCalled)
	}
}

func TestSendImage(t *testing.T) {
	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		image := []byte("FAKEIMAGE")
		SendImage(w, "image.png", "image/png", image)
	}))
	defer ts.Close()

	isCallbackCalled := false
	err := Request("GET", ts.URL, nil, func(res *http.Response) {
		// Make sure code reaches the callback
		isCallbackCalled = true

		// Make sure status code is as expected
		expectedStatus := 200
		if res.StatusCode != expectedStatus {
			t.Errorf("Expected %v but got %v", expectedStatus, res.StatusCode)
		}

		// Read image
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected no error, instead got %v", err.Error())
		}
		defer res.Body.Close()

		// Make sure the []byte of an image is being send
		expected := "FAKEIMAGE"
		actual := string(data)
		if expected != actual {
			t.Errorf("Expected %v but got %v", expected, actual)
		}

		// Make sure Content-Type is as expected
		expectedContentType := "image/png"
		actualContentType := res.Header.Get("Content-Type")
		if expectedContentType != actualContentType {
			t.Errorf("Expected %v but got %v", expectedContentType, actualContentType)
		}

		// Make sure Content-Disposition is as expected
		expectedContentDisposition := "inline; filename=image.png"
		actualContentDisposition := res.Header.Get("Content-Disposition")
		if expectedContentDisposition != actualContentDisposition {
			t.Errorf("Expected %v but got %v", expectedContentDisposition, actualContentDisposition)
		}
	})

	// We expect no error.
	if err != nil {
		t.Errorf("Expected no error, instead got %v", err.Error())
	}

	// To make sure that the callback is being called.
	expectedCb := true
	if isCallbackCalled != expectedCb {
		t.Errorf("Expected %v but got %v", expectedCb, isCallbackCalled)
	}
}
