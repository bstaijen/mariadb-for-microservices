package main

import (
	"net/http"
	"strconv"

	"github.com/bstaijen/mariadb-for-microservices/authentication-service/app/http/routes"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/config"
	"github.com/bstaijen/mariadb-for-microservices/authentication-service/database"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"

	// go-sql-driver/mysql is needed for the database connection
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/joho/godotenv/autoload"

	log "github.com/Sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	// Get config
	cnf := config.LoadConfig()

	// Get database
	connection, err := db.OpenConnection(cnf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConnection(connection)

	// Set the REST API routes
	routes := routes.InitRoutes(connection, cnf)
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(routes)

	// Start and listen on port in cbf.Port
	log.Info("Starting server on port " + strconv.Itoa(cnf.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cnf.Port), n))
}
