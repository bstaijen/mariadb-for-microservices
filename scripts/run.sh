#!/bin/bash
set -e

# Start running containers in the cluster and make 
# sure that registrator is running on each of them
docker-compose up -d --force-recreate
docker-compose scale registrator=5