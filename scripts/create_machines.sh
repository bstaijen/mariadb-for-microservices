#!/bin/bash
set -e

# Script is borrowed from:
# https://github.com/raykrueger/docker-swarm-mode-scripts/blob/master/create-swarm

echo "---Create manager1"
docker-machine create \
    --driver virtualbox \
    --virtualbox-memory 512 \
    manager1;

manager_ip=$(docker-machine ip manager1)

echo "---Swarm Init"
docker-machine ssh manager1 docker swarm init --listen-addr ${manager_ip} --advertise-addr ${manager_ip}

printf "\n---Get Tokens\n"
worker_token=$(docker-machine ssh manager1 docker swarm join-token -q worker)
echo ${worker_token}

for n in {1..4} ; do
	printf "\n---Create worker${n}\n"
	docker-machine create --driver virtualbox --virtualbox-memory 512 worker${n}
	ip=$(docker-machine ip worker${n})
	echo "--- Swarm Worker Join"
	docker-machine ssh worker${n} docker swarm join --listen-addr ${ip} --advertise-addr ${ip} --token ${worker_token} ${manager_ip}:2377
done