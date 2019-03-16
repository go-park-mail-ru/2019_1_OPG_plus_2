#!/usr/bin/env bash

docker stop $(cat ~/docker.log)
docker pr -v $(cat ~/docker.log)
docker image prune
rm ~/docker.log
echo
echo "Images left (just for checking whether no trash containers and images left)"
docker images
echo
docker ps -a