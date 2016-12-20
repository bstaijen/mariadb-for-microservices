package routes

import (
	"github.com/bstaijen/mariadb-for-microservices/comment-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/shared/util/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// InitRoutes instantiates a new gorilla/mux router
func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = setRESTRoutes(router)
	router = setIPCRoutes(router)
	return router
}

// setRESTRoutes specifies all public routes for the comment service
func setRESTRoutes(router *mux.Router) *mux.Router {
	// Subrouter /comments
	comments := router.PathPrefix("/comments").Subrouter()
	comments.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	// Create a comment /comments
	comments.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.CreateHandler),
	))
	return router
}

// Inter-Process Communication routes specifies the routes for internal communication
func setIPCRoutes(router *mux.Router) *mux.Router {
	// IPC subrouter /ipc
	ipc := router.PathPrefix("/ipc").Subrouter()

	// get last 10 comments /ipc/getLast10
	ipc.Handle("/getLast10", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.GetLastTenHandler),
	)).Methods("GET")
	return router
}
