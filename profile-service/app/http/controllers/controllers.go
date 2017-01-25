package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
	"github.com/gorilla/mux"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"

	log "github.com/Sirupsen/logrus"
)

// CreateUserHandler creates a new user in the database. Password is saved as a hash.
func CreateUserHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := &models.User{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("Bad json"))
			return
		}

		if err := user.Validate(); err == nil {
			if err := user.ValidatePassword(); err == nil {

				hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
				user.Hash = string(hash)

				createdID, err := db.CreateUser(connection, user)

				if err != nil {
					util.SendBadRequest(w, err)
					return
				}
				createdUser, _ := db.GetUserByID(connection, createdID)
				util.SendOK(w, createdUser)

			} else {
				util.SendBadRequest(w, err)
			}
		} else {
			util.SendBadRequest(w, err)
		}
	})
}

// DeleteUserHandler removes a user from the database. User can only deletes it's own record.
func DeleteUserHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			util.SendBadRequest(w, errors.New("token is mandatory"))
			return
		}

		user := &models.User{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("Bad json"))
			return
		}

		secretKey := cnf.SecretKey
		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64)

		if int(ID) != user.ID {
			util.SendBadRequest(w, errors.New("you can only delete your own user object"))
			return
		}

		db.DeleteUser(connection, user)
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}
		util.SendOK(w, string(""))

	})
}

// UpdateUserHandler updates an user based on it's user ID. User is only allowed to update it's own record. Verification is being done based on the JWT in the request.
func UpdateUserHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(string("token is mandatory")))
			return
		}

		user := &models.User{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("bad json"))
			return
		}

		secretKey := cnf.SecretKey
		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64) // gets the ID

		if int(ID) != user.ID {
			util.SendBadRequest(w, errors.New("you can only change your own user object"))
			return
		}

		if err := user.Validate(); err == nil {
			if err := user.ValidatePassword(); err == nil {

				db.UpdateUser(connection, user)

				util.SendOK(w, user)

			} else {
				util.SendBadRequest(w, err)
			}
		} else {
			util.SendBadRequest(w, err)
		}
	})
}

// UserByIndexHandler retrieves an user from the database based on its id. This handler expects the id being passed in the route variable in the current request.
func UserByIndexHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		vars := mux.Vars(r)
		strID := vars["id"]

		id, err := strconv.Atoi(strID)
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}

		user, err := db.GetUserByID(connection, id)

		if err != nil {
			util.SendBadRequest(w, err)
			return
		}
		util.SendOK(w, user)
	})
}

// GetUsernamesHandler is a handler for collecting usernames coresponding to user ID's. The request expects a json object in the following format: { "requests":[{"id":1} ,{"id":2},{"id":3}, {"id":4} ]}.
func GetUsernamesHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result, err := bodyToArrayWithIDs(r)

		if err != nil {
			util.SendError(w, err)
			return
		}

		users, err := db.GetUsernames(connection, result)
		if err != nil {
			util.SendError(w, err)
			return
		}

		type Resp struct {
			Usernames []*sharedModels.GetUsernamesResponse `json:"usernames"`
		}
		util.SendOK(w, &Resp{Usernames: users})
	})
}

// Converts a json object to a list of ID's. Expects JSON to be in the following format: {"requests":[{"id":1},{"id":2},{"id":3},{"id":4} ]}
func bodyToArrayWithIDs(req *http.Request) ([]*sharedModels.GetUsernamesRequest, error) {
	type Collection struct {
		Objects []*sharedModels.GetUsernamesRequest `json:"requests"`
	}
	col := &Collection{}
	col.Objects = make([]*sharedModels.GetUsernamesRequest, 0)

	err := util.RequestToJSON(req, &col)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return col.Objects, nil
}
