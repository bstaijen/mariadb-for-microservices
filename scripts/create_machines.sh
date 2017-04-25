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
manager_token=$(docker-machine ssh manager1 docker swarm join-token -q manager)
worker_token=$(docker-machine ssh manager1 docker swarm join-token -q worker)
echo ${manager_token}
echo ${worker_token}

for n in {1..4} ; do
	printf "\n---Create worker${n}\n"
	docker-machine create --driver virtualbox --virtualbox-memory 512 worker${n}
	ip=$(docker-machine ip worker${n})
	echo "--- Swarm Worker Join"
	docker-machine ssh worker${n} docker swarm join --listen-addr ${ip} --advertise-addr ${ip} --token ${worker_token} ${manager_ip}:2377
done

printf "\n---Launching Visualizer\n"
docker-machine ssh manager1 docker run -it -d -p 8080:8080 -e HOST=${manager_ip} -e PORT=8080 -v /var/run/docker.sock:/var/run/docker.sock manomarks/visualizer


printf "\n---Launching ETCD\n"
eval $(docker-machine env manager1)
docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
 --name etcd quay.io/coreos/etcd etcd\
 -name etcd0 \
 -advertise-client-urls http://${manager_ip}:2379,http://${manager_ip}:4001 \
 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
 -initial-advertise-peer-urls http://${manager_ip}:2380 \
 -listen-peer-urls http://0.0.0.0:2380 \
 -initial-cluster-token etcd-cluster-1 \
 -initial-cluster etcd0=http://${manager_ip}:2380 \
 -initial-cluster-state new

docker network create --driver overlay --opt encrypted my-network

printf "\n\n------------------------------------\n"
echo "To visualize your cluster..."
echo "Open a browser to http://${manager_ip}:8080/"
echo "To connect to your cluster..."
echo 'eval $(docker-machine env manager1)'
