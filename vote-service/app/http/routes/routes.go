package routes

import (
	"github.com/bstaijen/mariadb-for-microservices/vote-service/app/http/controllers"

	"github.com/bstaijen/mariadb-for-microservices/shared/util/middleware"
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
	votes := router.PathPrefix("/votes").Subrouter()
	votes.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))
	votes.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.CreateHandler),
	))
	return router
}

// Inter-Process Communication routes
func setIPCRoutes(router *mux.Router) *mux.Router {
	ipc := router.PathPrefix("/ipc").Subrouter()
	ipc.Handle("/toprated", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.GetTopRatedHandler),
	)).Methods("GET")
	ipc.Handle("/hot", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.GetHotHandler),
	)).Methods("GET")
	ipc.Handle("/count", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.GetVoteCountHandler),
	))
	ipc.Handle("/voted", negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		negroni.HandlerFunc(controllers.HasVotedHandler),
	))
	return router
}
