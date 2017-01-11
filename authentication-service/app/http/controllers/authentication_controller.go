package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/config"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler validates the user and returns a JWT token
func LoginHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	login := &Login{}

	err := util.RequestToJSON(r, login)
	if err != nil {
		util.SendBadRequest(w, errors.New("Bad json"))
		return
	}

	// Check if there's atleast some value
	if len(login.Username) < 1 || len(login.Password) < 1 {
		util.SendErrorMessage(w, "Please provide username and password in the body")
		return
	}

	// authenticate the username password combination
	usr, err := authenticate(login.Username, login.Password)
	if err != nil {
		util.SendBadRequest(w, err)
		return
	}

	// create JWT object with claims
	expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.ID,
		"iat": time.Now().Unix(),
		"exp": expiration,
	})

	// Load secret key from config and generate a signed token
	secretKey := config.LoadConfig().SecretKey
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		util.SendError(w, err)
		return
	}

	// Retrieve the user from the database
	database := db.InitMariaDB()
	databaseUser, _ := database.GetUserByUsername(login.Username)

	// Send the token and user back
	data := &models.Token{
		Token:     tokenString,
		ExpiresOn: strconv.Itoa(int(expiration)),
		User:      databaseUser,
	}
	util.SendOK(w, data)
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	util.SendErrorMessage(w, "Not implemented!")
}

// authenticate user by checking username and password in database
func authenticate(username string, password string) (*models.User, error) {

	database := db.InitMariaDB()

	databaseUser, _ := database.GetUserByUsername(username)

	if bcrypt.CompareHashAndPassword([]byte(databaseUser.Password), []byte(password)) == nil {
		return &databaseUser, nil
	}
	return &models.User{}, ErrInvalidCredentials
}

// ErrInvalidCredentials error
var ErrInvalidCredentials = errors.New("Invalid credentials")
