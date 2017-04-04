#!/bin/bash
set -e

# Docker Swarm Setup

# Create consul
docker-machine create \
    --driver virtualbox \
    --virtualbox-memory 512 consul

eval "$(docker-machine env consul)"

docker run --restart=always -d \
    -p "8500:8500" \
    -h "consul" progrium/consul -server -bootstrap

# Create master node
docker-machine create \
    --driver virtualbox \
    --virtualbox-memory 512 \
   --swarm --swarm-master \
    --swarm-discovery="consul://$(docker-machine ip consul):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip consul):8500" \
    --engine-opt="cluster-advertise=eth1:2376" master

# Create 4 nodes
for N in 0 1; do 
    docker-machine create \
        --driver virtualbox \
        --virtualbox-memory 512 \
        --swarm \
        --swarm-discovery="consul://$(docker-machine ip consul):8500" \
        --engine-opt="cluster-store=consul://$(docker-machine ip consul):8500" \
        --engine-opt="cluster-advertise=eth1:2376" node$N;
done