eval $(docker-machine env dbmaster)
docker rmi $(docker images -q) -f

eval $(docker-machine env master)
docker rmi $(docker images -q) -f

eval $(docker-machine env node1)
docker rmi $(docker images -q) -f

eval $(docker-machine env node2)
docker rmi $(docker images -q) -f

eval $(docker-machine env node3)
docker rmi $(docker images -q) -f

eval $(docker-machine env node4)
docker rmi $(docker images -q) -f

eval $(docker-machine env node5)
docker rmi $(docker images -q) -f