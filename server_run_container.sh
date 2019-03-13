#!/usr/bin/env bash

if [[ "$#" -eq 1 ]]
then
    port=$1
    else
    port=8002
    fi

echo The published port of the container is ${port}

docker build --tag colors:back -f ServerDockerfile --rm
docker run -d -e PORT=${port} --name colors-back --publish ${port}:8002 colors:back >> docker.log
