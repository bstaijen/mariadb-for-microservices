package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/config"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler validates the user and returns a JWT token
func LoginHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
			util.SendBadRequest(w, errors.New("Please provide username and password in the body"))
			return
		}

		// authenticate the username password combination
		usr, err := authenticate(connection, login.Username, login.Password)
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
		secretKey := cnf.SecretKey
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			util.SendError(w, err)
			return
		}

		// Retrieve the user from the database
		databaseUser, _ := db.GetUserByUsername(connection, login.Username)

		// Send the token and user back
		data := &models.Token{
			Token:     tokenString,
			ExpiresOn: strconv.Itoa(int(expiration)),
			User:      databaseUser,
		}
		util.SendOK(w, data)
	})
}

// authenticate user by checking username and password in database
func authenticate(connection *sql.DB, username string, password string) (*models.User, error) {
	databaseUser, _ := db.GetUserByUsername(connection, username)
	if bcrypt.CompareHashAndPassword([]byte(databaseUser.Password), []byte(password)) == nil {
		return &databaseUser, nil
	}
	return &models.User{}, ErrInvalidCredentials
}

// ErrInvalidCredentials error
var ErrInvalidCredentials = errors.New("Invalid credentials")
