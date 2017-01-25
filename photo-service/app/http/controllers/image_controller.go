package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"errors"

	"io/ioutil"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/models"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
)

func CreateHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
	PanicIfError(err)
	defer file.Close()

	image, err := ioutil.ReadAll(file)
	PanicIfError(err)

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
	database := db.InitMariaDB()
	err = database.InsertPhoto(img)
	if err != nil {
		util.SendBadRequest(w, err)
		return
	}

	util.SendOK(w, string("Success"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	file := vars["file"]
	// TODO : what if file not exist?
	database := db.InitMariaDB()
	photo, err := database.GetPhotoByFilename(file)
	if err != nil {
		util.SendError(w, err)
		return
	}
	util.SendImage(w, photo.Filename, photo.ContentType, photo.Image)
}

func ListByUserIDHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	strID := vars["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		util.SendErrorMessage(w, "id must be integer")
		return
	}

	database := db.InitMariaDB()
	photos, err := database.ListImagesByUserID(id)
	if err != nil {
		util.SendError(w, err)
		return
	}
	util.SendOK(w, photos)
}

func randomFileName() string {
	rand.Seed(time.Now().UTC().UnixNano())
	var res string
	for i := 0; i <= 16; i++ {
		res += strconv.Itoa(rand.Intn(10))
	}
	return res
}

func PanicIfError(err error) {
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		panic(err)
	}
}
