package routes

import (
	"database/sql"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/middleware"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"

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

	// Subrouter /users
	users := router.PathPrefix("/users").Subrouter()

	// TODO :  https://github.com/gorilla/handlers/blob/master/cors.go#L140
	users.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONSHandler),
	))

	// Update user /users
	users.Methods("PUT").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		middleware.RequireTokenAuthenticationHandler(cnf),
		controllers.UpdateUserHandler(db, cnf),
	))

	// Delete User /users
	users.Methods("DELETE").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		middleware.RequireTokenAuthenticationHandler(cnf),
		controllers.DeleteUserHandler(db, cnf),
	))

	// Create user /sers
	users.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.CreateUserHandler(db),
	))

	// Get one user /user/{id}
	oneUser := router.PathPrefix("/user/{id}").Subrouter()
	oneUser.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		middleware.RequireTokenAuthenticationHandler(cnf),
		controllers.UserByIndexHandler(db),
	))

	return router
}

// Inter-Process Communication routes
func setIPCRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {
	// IPC subrouter /ipc
	ipc := router.PathPrefix("/ipc").Subrouter()

	// get usernames /ipc/usernames
	ipc.Handle("/usernames", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.GetUsernamesHandler(db),
	)).Methods("GET")

	return router
}
