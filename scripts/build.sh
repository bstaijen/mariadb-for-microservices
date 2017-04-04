#!/bin/bash
set -e

# Build the microservice images
docker build -t webserver ../webserver
docker build -t authenticationsvc ../authentication-service
docker build -t photosvc ../photo-service
docker build -t votesvc ../vote-service
docker build -t commentsvc ../comment-service
docker build -t profilesvc ../profile-service
docker build -t vizualizersvc ../vizualizer
docker build -t db ../database