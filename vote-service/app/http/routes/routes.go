package routes

import (
	"database/sql"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/vote-service/config"

	"github.com/bstaijen/mariadb-for-microservices/shared/util/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// InitRoutes initializes the REST and IPC routes for this service.
func InitRoutes(db *sql.DB, cnf config.Config) *mux.Router {
	router := mux.NewRouter()
	router = setRESTRoutes(db, cnf, router)
	router = setIPCRoutes(db, cnf, router)
	return router
}

func setRESTRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {
	votes := router.PathPrefix("/votes").Subrouter()
	votes.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))
	votes.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.CreateHandler(db, cnf),
	))
	return router
}

// Inter-Process Communication routes
func setIPCRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {
	ipc := router.PathPrefix("/ipc").Subrouter()
	ipc.Handle("/toprated", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.GetTopRatedHandler(db),
	)).Methods("GET")
	ipc.Handle("/hot", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.GetHotHandler(db),
	)).Methods("GET")
	ipc.Handle("/count", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.GetVoteCountHandler(db),
	))
	ipc.Handle("/voted", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.HasVotedHandler(db),
	))
	return router
}
