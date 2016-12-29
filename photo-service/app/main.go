package main

import (
	"log"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/app/http/routes"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"

	"strconv"

	"github.com/urfave/negroni"

	// go-sql-driver/mysql is needed for the database connection
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// Get config
	cnf := config.LoadConfig()

	// Set the REST API routes
	router := routes.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	// Start and listen on port in cbf.Port
	log.Println("Starting server on port " + strconv.Itoa(cnf.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cnf.Port), n))
}
