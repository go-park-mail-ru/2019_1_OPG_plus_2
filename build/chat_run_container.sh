#!/usr/bin/env bash

if [[ "$#" -eq 1 ]]
then
    port=$1
    else
    port=8003
    fi

echo [SERVER-CHAT] The published port of the container is ${port}

docker build --rm -t colors-back-chat -f chat.Dockerfile .. &&\
docker run -d --network=opg-net -p ${port}:8003 --name colors-back-chat colors-back-chat >> ~/docker.log &&\
echo [SERVER-CHAT] Server is now running at: ${port}