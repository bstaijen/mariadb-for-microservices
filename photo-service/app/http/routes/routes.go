package routes

import (
	"database/sql"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
	"github.com/bstaijen/mariadb-for-microservices/shared/util/middleware"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// InitRoutes instantiates a new gorilla/mux router
func InitRoutes(db *sql.DB, cnf config.Config) *mux.Router {
	router := mux.NewRouter()
	router = setPhotoRoutes(db, cnf, router)
	router = setIPCRoutes(db, cnf, router)
	return router
}

// setPhotoRoutes specifies all routes for the authentication service
func setPhotoRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {

	// Subrouter /image
	image := router.PathPrefix("/image").Subrouter()

	// Options
	image.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	image.Handle("/{id}/delete", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.DeletePhotoHandler(db, cnf),
	)).Methods("POST")

	// Add image for user /image/{id}
	image.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(middleware.RequireTokenAuthenticationHandler(cnf.SecretKey)),
		controllers.CreateHandler(db),
	)).Methods("POST")

	image.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	)).Methods("OPTIONS")

	// Image for user /image/{id}/list
	image.Handle("/{id}/list", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.ListByUserIDHandler(db, cnf),
	)).Methods("GET")

	// Incoming Timeline /image/list
	image.Handle("/list", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.IncomingHandler(db, cnf),
	)).Methods("GET")

	// Top Rated Timeline /image/toprated
	image.Handle("/toprated", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.TopRatedHandler(db, cnf),
	)).Methods("GET")

	// Hot Timeline /image/hot
	image.Handle("/hot", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.HotHandler(db, cnf),
	)).Methods("GET")

	// Subrouter /images/{file}
	images := router.PathPrefix("/images/{file}").Subrouter()

	// Retrieve single image /images/{file}
	images.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.IndexHandler(db),
	))

	return router
}

func setIPCRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {
	// Subrouter /ipc
	image := router.PathPrefix("/ipc").Subrouter()

	// Options
	image.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	// Get photos
	image.Handle("/getPhotos", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.IPCGetPhotos(db, cnf),
	)).Methods("GET")

	return router
}
