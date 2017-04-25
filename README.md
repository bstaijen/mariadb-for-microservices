[![Build Status](https://travis-ci.org/bstaijen/mariadb-for-microservices.svg?branch=master)](https://travis-ci.org/bstaijen/mariadb-for-microservices) [![Go Report Card](https://goreportcard.com/badge/github.com/bstaijen/mariadb-for-microservices)](https://goreportcard.com/report/github.com/bstaijen/mariadb-for-microservices)
Work in progress..

# MariaDB for Microservices
MariaDB for Microservices is a working example on how to create an application using the microservice architectural approach and the MariaDB Server.

## How does it work?
This project consists out of 5 services, a shared library, 1 database and a webserver (also serves as proxy). We use docker for to deploy and run each service. 

## Requirements
- Docker (version 17.03.1)

## Usage
To run the application we first have to set-up a cluster of machines, after that we have to configure our database, we end with bootstrapping the application.
Create a cluster of 5 machines with the following command: `./create_machines.sh`. This script also install a visualizer and etcd for service discovery.

Now pass the location of the etcd service to the database configuration. To do that you have to edit the docker-compose-stacks.yml file and add the IP and port to the “DISCOVERY_SERVICE” environment variable of the database. The configuration should be on line 126. You can get the IP by using the command `docker-machine ip manager-1`.

The last thing to do is to deploy the application, use the command `docker stack deploy --compose-file docker-compose-stacks.yml demo`

# Feedback & Issues
- Feel free to report bugs or suggestions through the Github issues page.

# Authors
- Bjorge Staijen