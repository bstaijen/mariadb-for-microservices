package util

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"errors"

	"github.com/bstaijen/mariadb-for-microservices/shared/models"
)

// SendOKMessage send http response with header OK and a message.
func SendOKMessage(w http.ResponseWriter, message string) {
	SendOK(w, &models.Error{Message: message})
}

// SendOK send http response with header OK and a body interface.
func SendOK(w http.ResponseWriter, data interface{}) {
	result, err := json.Marshal(data)
	if err != nil {
		SendBadRequest(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

// SendErrorMessage send a http response an error message
func SendErrorMessage(w http.ResponseWriter, message string) {
	SendError(w, errors.New(message))
}

// SendError send a http response with status code 400 with the error parameter as error message.
func SendError(w http.ResponseWriter, err error) {
	debug.PrintStack()
	SendBadRequest(w, err)
}

// SendBadRequest writes a Bad Request to the ResponseWrite
func SendBadRequest(w http.ResponseWriter, err error) {
	e := &models.Error{Message: err.Error()}
	var errJSON, _ = json.Marshal(e)

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(errJSON)))
}

// SendImage send a http response with a write a image to the client
func SendImage(w http.ResponseWriter, filename string, contentType string, image []byte) {
	w.Header().Set("Content-Disposition", "inline; filename="+filename)
	w.Header().Set("Content-Type", contentType)
	w.Write(image)
}
