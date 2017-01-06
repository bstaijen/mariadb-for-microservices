package main

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/app/http/routes"
	"github.com/bstaijen/mariadb-for-microservices/vote-service/config"
	negronilogrus "github.com/meatballhat/negroni-logrus"

	"github.com/urfave/negroni"

	// go-sql-driver/mysql is needed for the database connection
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Get config
	cnf := config.LoadConfig()

	// Set routes
	routes := routes.InitRoutes()
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(routes)

	// Start and listen on port in cbf.Port
	log.Info("Starting server on port " + strconv.Itoa(cnf.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cnf.Port), n))
}
