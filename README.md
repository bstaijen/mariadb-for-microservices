[![Build Status](https://travis-ci.org/bstaijen/mariadb-for-microservices.svg?branch=master)](https://travis-ci.org/bstaijen/mariadb-for-microservices) [![Go Report Card](https://goreportcard.com/badge/github.com/bstaijen/mariadb-for-microservices)](https://goreportcard.com/report/github.com/bstaijen/mariadb-for-microservices)
Work in progress..

# MariaDB for Microservices
MariaDB for Microservices is a working example on how to create an application using the microservice architectural approach and the MariaDB Server.

## How does it work?
This project consists out of 5 services, a shared library, 1 database and a webserver (also serves as proxy). We use docker for to deploy and run each service. 

## Requirements
- Docker (version 1.12.5)
- Docker Swarm
- Docker Compose (version 1.9.0)

## Usage
The quickest way to run the application is by running the following commands:
`docker-compose up -d --force-recreate`
`docker-compose scale registrator=7` The database is dependent on registrator. Set the number of registrators to the number of machines in your cluster.
Best way to execute this is to create a shell script. 


# Feedback & Issues
- Feel free to report bugs or suggestions through the Github issues page.

# Authors
- Bjorge Staijen

# Next steps
Check out the [docs directory](docs) for more docs.