package routes

import (
	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/shared/util/middleware"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetRoutes(router)
	return router
}

func SetRoutes(router *mux.Router) *mux.Router {

	// Subroutr /image
	image := router.PathPrefix("/image").Subrouter()

	// Options
	image.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	// Add image for user /image/{id}
	image.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.CreateHandler),
	)).Methods("POST")

	image.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	)).Methods("OPTIONS")

	// Image for user /image/{id}/list
	image.Handle("/{id}/list", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.ListByUserIDHandler),
	)).Methods("GET")

	// Incoming Timeline /image/list
	image.Handle("/list", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.IncomingHandler),
	)).Methods("GET")

	// Top Rated Timeline /image/toprated
	image.Handle("/toprated", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.TopRatedHandler),
	)).Methods("GET")

	// Hot Timeline /image/hot
	image.Handle("/hot", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.HotHandler),
	)).Methods("GET")

	// Subrouter /images/{file}
	images := router.PathPrefix("/images/{file}").Subrouter()

	// Retrieve single image /images/{file}
	images.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.IndexHandler),
	))

	return router
}
