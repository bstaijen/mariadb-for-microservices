package routes

import (
	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = setRESTRoutes(router)
	router = setIPCRoutes(router)
	return router
}

func setRESTRoutes(router *mux.Router) *mux.Router {

	// Subrouter /users
	users := router.PathPrefix("/users").Subrouter()
	users.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	// Update user /users
	users.Methods("PUT").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(middleware.RequireTokenAuthenticationController),
		negroni.HandlerFunc(controllers.UpdateUserController),
	))

	// Delete User /users
	users.Methods("DELETE").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(middleware.RequireTokenAuthenticationController),
		negroni.HandlerFunc(controllers.DeleteUserController),
	))

	// Create user /sers
	users.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.CreateUserController),
	))

	// Get one user /user/{id}
	oneUser := router.PathPrefix("/user/{id}").Subrouter()
	oneUser.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(middleware.RequireTokenAuthenticationController),
		negroni.HandlerFunc(controllers.UserIndexController),
	))

	return router
}

// Inter-Process Communication routes
func setIPCRoutes(router *mux.Router) *mux.Router {
	// IPC subrouter /ipc
	ipc := router.PathPrefix("/ipc").Subrouter()

	// get usernames /ipc/usernames
	ipc.Handle("/usernames", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.GetUsernames),
	)).Methods("GET")

	return router
}
