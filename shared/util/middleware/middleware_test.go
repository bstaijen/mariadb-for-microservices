package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

	// Test if the code gets executed when there's a next arguments passed
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

// Test getting and validating token in the query parameter.
func TestRequireTokenInURLParam(t *testing.T) {

	// Create JWT object with claims
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	tokenString, err := token.SignedString([]byte("ABCDEF"))
	if err != nil {
		t.Error(err)
		return
	}

	// Make rquest
	req, err := http.NewRequest("GET", "http://localhost/test?token="+tokenString, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Invoke middleware
	handler := RequireTokenAuthenticationHandler("ABCDEF")
	handler(res, req, nil)

	// Test results:

	// In case of no token
	if res.Result().StatusCode == 400 {
		t.Errorf("Expected statuscode to be empty but got %v. Bad Request.", res.Result().StatusCode)
	}

	// In case of bad token
	if res.Result().StatusCode == 404 {
		t.Errorf("Expected statuscode to be empty but got %v. Unauthorized.", res.Result().StatusCode)
	}

	// Nothing supposed to happen if everything goes right.
}

// Test getting and validating token from request header.
func TestRequireTokenInHeader(t *testing.T) {
	// Create JWT object with claims
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	tokenString, err := token.SignedString([]byte("ABCDEF"))
	if err != nil {
		t.Error(err)
		return
	}

	// Make rquest
	req, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", tokenString)
	res := httptest.NewRecorder()

	// Invoke middleware
	handler := RequireTokenAuthenticationHandler("ABCDEF")
	handler(res, req, nil)

	// Test results:

	// In case of no token
	if res.Result().StatusCode == 400 {
		t.Errorf("Expected statuscode to be empty but got %v. Bad Request.", res.Result().StatusCode)
	}

	// In case of bad token
	if res.Result().StatusCode == 404 {
		t.Errorf("Expected statuscode to be empty but got %v. Unauthorized.", res.Result().StatusCode)
	}

	// Nothing supposed to happen if everything goes right.
}

func TestNoTokenProvided(t *testing.T) {
	// Make rquest
	req, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Invoke middleware
	handler := RequireTokenAuthenticationHandler("")
	handler(res, req, nil)

	if res.Result().StatusCode != 400 {
		t.Errorf("Expected statuscode to be 400 but got %v.", res.Result().StatusCode)
	}
}

func TestBadToken(t *testing.T) {
	// Make rquest
	req, err := http.NewRequest("GET", "http://localhost/test?token=token", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Invoke middleware
	handler := RequireTokenAuthenticationHandler("ThisIsNotAGoodSecretKey")
	handler(res, req, nil)

	if res.Result().StatusCode != 400 {
		t.Errorf("Expected statuscode to be 400 but got %v.", res.Result().StatusCode)
	}
}

func TestTokenOK(t *testing.T) {
	// Create JWT object with claims
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	tokenString, err := token.SignedString([]byte("ABCDEF"))
	if err != nil {
		t.Error(err)
		return
	}

	// Make rquest
	req, err := http.NewRequest("GET", "http://localhost/test?token="+tokenString, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Invoke middleware
	handler := RequireTokenAuthenticationHandler("ABCDEF")
	handler(res, req, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test", "value")
	})

	if res.Result().StatusCode != 200 {
		t.Errorf("Expected statuscode to be 200 but got %v.", res.Result().StatusCode)
	}
}

func TestExpiredToken(t *testing.T) {
	// Create JWT object with claims
	expiration := (time.Now().Unix() - 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Generate a signed token
	tokenString, err := token.SignedString([]byte("ABCDEF"))
	if err != nil {
		t.Error(err)
		return
	}

	// Make rquest
	req, err := http.NewRequest("GET", "http://localhost/test?token="+tokenString, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	// Invoke middleware
	handler := RequireTokenAuthenticationHandler("ABCDEF")
	handler(res, req, nil)

	if res.Result().StatusCode != 400 {
		t.Errorf("Expected statuscode to be 400 but got %v.", res.Result().StatusCode)
	}
}
