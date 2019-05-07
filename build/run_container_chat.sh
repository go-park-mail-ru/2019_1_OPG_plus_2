#!/usr/bin/env bash

port=8003
while getopts 'p:' flag; do
  case ${flag} in
    p) port=${OPTARG};;
  esac
done

echo [SERVER-CHAT] The published port of the container is ${port}

docker build --rm -t colors-back-chat -f dockerfiles/chat.Dockerfile .. &&\
docker run -d \
           --network=opg-net \
           -p ${port}:8003 \
           --name colors-back-chat \
           colors-back-chat >> ~/docker.log &&\
echo [SERVER-CHAT] Server is now running at: ${port}