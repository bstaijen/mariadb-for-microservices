package routes

import (
	"database/sql"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
	"github.com/bstaijen/mariadb-for-microservices/shared/util/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// InitRoutes instantiates a new gorilla/mux router
func InitRoutes(db *sql.DB, cnf config.Config) *mux.Router {
	router := mux.NewRouter()
	router = setRESTRoutes(db, cnf, router)
	router = setIPCRoutes(db, cnf, router)
	return router
}

// setRESTRoutes specifies all public routes for the comment service
func setRESTRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {
	// Subrouter /comments
	comments := router.PathPrefix("/comments").Subrouter()
	comments.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	comments.Handle("/fromuser", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.ListCommentsFromUser(db, cnf),
	)).Methods("GET")

	// Create a comment /comments
	comments.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.CreateHandler(db, cnf),
	))
	comments.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.ListCommentsHandler(db, cnf),
	))

	return router
}

// Inter-Process Communication routes specifies the routes for internal communication
func setIPCRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {
	// IPC subrouter /ipc
	ipc := router.PathPrefix("/ipc").Subrouter()

	// get last 10 comments /ipc/getLast10
	ipc.Handle("/getLast10", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.GetLastTenHandler(db, cnf),
	)).Methods("GET")

	// get the number of comments of a photo /ipc/getCount
	ipc.Handle("/getCount", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.GetCommentCountHandler(db, cnf),
	)).Methods("GET")

	return router
}
