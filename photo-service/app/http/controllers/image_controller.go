package controllers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"errors"

	"io/ioutil"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

// CreateHandler create a photo object and store it in the database.
func CreateHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// Get title
		var title = r.URL.Query().Get("title")
		if len(title) < 1 {
			util.SendBadRequest(w, errors.New("Title is mandatory"))
		}

		// Get userID
		vars := mux.Vars(r)
		strID := vars["id"]
		id, _ := strconv.Atoi(strID)

		// Read file
		file, fileheader, err := r.FormFile("file")
		if err != nil {
			util.SendError(w, err)
			return
		}
		defer file.Close()

		image, err := ioutil.ReadAll(file)
		if err != nil {
			util.SendError(w, err)
			return
		}

		// Get extensions, filename and contenttype
		extension := strings.Split(fileheader.Filename, ".")[1]
		filename := fmt.Sprintf("%v.%v", randomFileName(), extension)
		contentType := fileheader.Header.Get("Content-Type")

		// Create model
		img := &models.CreatePhoto{
			UserID:      id,
			Filename:    filename,
			Title:       title,
			ContentType: contentType,
			Image:       image,
		}

		// Save
		err = db.InsertPhoto(connection, img)
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}

		util.SendOK(w, string("Success"))
	})
}

// IndexHandler serves a photo indentiefied by filename
func IndexHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		file := vars["file"]
		// TODO : what if file not exist?

		photo, err := db.GetPhotoByFilename(connection, file)
		if err != nil {
			util.SendError(w, err)
			return
		}
		util.SendImage(w, photo.Filename, photo.ContentType, photo.Image)
	})
}

// ListByUserIDHandler list all photos owned by an user.
func ListByUserIDHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		strID := vars["id"]
		id, err := strconv.Atoi(strID)
		if err != nil {
			util.SendErrorMessage(w, "id must be integer")
			return
		}

		photos, err := db.ListImagesByUserID(connection, id)
		if err != nil {
			util.SendError(w, err)
			return
		}

		photos = findResources(cnf, photos, id, true, true, true)

		util.SendOK(w, photos)
	})
}

func GetPhotoByID(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		userID, err := getUserIDFromRequest(cnf, r)
		if err != nil {
			util.SendError(w, err)
			return
		}

		vars := mux.Vars(r)
		strID := vars["id"]
		id, err := strconv.Atoi(strID)
		if err != nil {
			util.SendErrorMessage(w, fmt.Sprintf("id must be integer, instead got %v", id))
			return
		}

		photo, err := db.GetPhotoById(connection, id)
		if err != nil {
			util.SendError(w, err)
			return
		}

		logrus.Info(photo.ID)

		photos := make([]*models.Photo, 0)
		photos = append(photos, photo)
		photos = findResources(cnf, photos, userID, true, true, true)

		util.SendOK(w, photos[0])
	})
}

func randomFileName() string {
	rand.Seed(time.Now().UTC().UnixNano())
	var res string
	for i := 0; i <= 16; i++ {
		res += strconv.Itoa(rand.Intn(10))
	}
	return res
}
