package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bstaijen/mariadb-for-microservices/vote-service/app/http/routes"
	"github.com/bstaijen/mariadb-for-microservices/vote-service/config"

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
	n.UseHandler(routes)

	// Start and listen on port in cbf.Port
	log.Println("Starting server on port " + strconv.Itoa(cnf.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cnf.Port), n))
}
