package main

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/bstaijen/mariadb-for-microservices/profile-service/app/http/routes"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/config"
	"github.com/bstaijen/mariadb-for-microservices/profile-service/database"
	negronilogrus "github.com/meatballhat/negroni-logrus"

	"github.com/urfave/negroni"

	// go-sql-driver/mysql is needed for the database connection
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	log.SetLevel(log.DebugLevel)

	// Get config
	cnf := config.LoadConfig()

	// Get database
	connection, _ := db.OpenConnection()
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
