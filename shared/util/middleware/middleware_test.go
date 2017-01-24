package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test if AccessControlHandler sets the Access-Control-Allow-Origin header.
func TestAccessControlHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	AccessControlHandler(res, req, nil)

	exp := "*"
	act := res.Header().Get("Access-Control-Allow-Origin")
	if exp != act {
		t.Fatalf("Expected %s got %s", exp, act)
	}
}

// Test if the next function is being executed.
func TestAccessControlHandlerWithNext(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Test if the code gets exectued when there's a next arguments passed
	AccessControlHandler(res, req, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test", "value")
	})

	exp := "*"
	act := res.Header().Get("Access-Control-Allow-Origin")
	if exp != act {
		t.Fatalf("Expected %s got %s", exp, act)
	}

	exp = "value"
	act = res.Header().Get("test")
	if exp != act {
		t.Fatalf("Expected %s got %s", exp, act)
	}
}

// Test if the Access-Control-Allow-Origin and Access-Control-Allow-Headers are being set by the AcceptOPTIONSHandler middleware
func TestAcceptOPTIONS(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	AcceptOPTIONS(res, req, nil)

	exp := "*"
	act := res.Header().Get("Access-Control-Allow-Origin")
	if exp != act {
		t.Fatalf("Expected %s got %s", exp, act)
	}

	exp = "Content-Type"
	act = res.Header().Get("Access-Control-Allow-Headers")
	if exp != act {
		t.Fatalf("Expected %s got %s", exp, act)
	}
}
