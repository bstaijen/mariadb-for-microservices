package routes

import (
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/http/controllers"
	"github.com/bstaijen/mariadb-for-microservices/shared/util/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// InitRoutes instantiates a new gorilla/mux router
func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = setAuthenticationRoutes(router)
	return router
}

// setAuthenticationRoutes specifies all routes for the authentication service
func setAuthenticationRoutes(router *mux.Router) *mux.Router {

	// Subrouter /token-auth
	tokenAUTH := router.PathPrefix("/token-auth").Subrouter()

	// User Login POST /token-auth
	tokenAUTH.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.LoginHandler),
	))

	// OPTIONS /token-auth
	tokenAUTH.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	// User Refresh token. GET /refresh-token-auth
	router.Handle("/refresh-token-auth",
		negroni.New(
			negroni.HandlerFunc(middleware.AccessControlHandler),
			negroni.HandlerFunc(controllers.RefreshTokenHandler),
		)).Methods("GET")
	return router
}
