#!/usr/bin/env bash

if [[ "$#" -eq 1 ]]
then
    port=$1
    else
    port=8002
    fi

echo [SERVER] The published port of the container is ${port}

docker build --rm -t colors-back -f server.Dockerfile .
docker run -d --rm --network=opg-net --name colors-back --publish ${port}:8002 colors-back >> docker.log

echo [SERVER] Server is now running at: ${port}