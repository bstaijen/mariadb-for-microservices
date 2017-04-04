#!/bin/bash
set -e

SWARM_MACHINE_NAME=master
CERT_PATH=~/.docker/machine/machines/$SWARM_MACHINE_NAME
export SWARM_CA=$(cat $CERT_PATH/ca.pem)
export SWARM_CERT=$(cat $CERT_PATH/cert.pem)
export SWARM_KEY=$(cat $CERT_PATH/key.pem)
export SWARM_HOST=tcp://$(docker-machine ip $SWARM_MACHINE_NAME):3376

# Start running containers in the cluster and make 
# sure that registrator is running on each of them
docker-compose -f docker-compose-registrator.yml up -d
docker-compose -f docker-compose-registrator.yml scale registrator=3