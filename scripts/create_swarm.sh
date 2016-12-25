#!/bin/bash

set -e

# Docker Swarm Setup
# Create machine
docker-machine create \
    --driver virtualbox \
    --virtualbox-memory 512 consul

eval "$(docker-machine env consul)"

docker run --restart=always -d \
    -p "8500:8500" \
    -h "consul" progrium/consul -server -bootstrap

# Create master
docker-machine create \
    --driver virtualbox \
    --virtualbox-memory 512 \
   --swarm --swarm-master \
    --swarm-discovery="consul://$(docker-machine ip consul):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip consul):8500" \
    --engine-opt="cluster-advertise=eth1:2376" master

# Create 5 
for N in 0 1 2 3 4 5; do 
    docker-machine create \
        --driver virtualbox \
        --virtualbox-memory 512 \
        --swarm \
        --swarm-discovery="consul://$(docker-machine ip consul):8500" \
        --engine-opt="cluster-store=consul://$(docker-machine ip consul):8500" \
        --engine-opt="cluster-advertise=eth1:2376" node$N;
done