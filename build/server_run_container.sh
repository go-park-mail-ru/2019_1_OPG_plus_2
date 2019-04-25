#!/usr/bin/env bash

if [[ "$#" -eq 1 ]]
then
    port=$1
    else
    port=8002
    fi

echo [SERVER-CORE] The published port of the container is ${port}

docker build --rm -t colors-back-core -f core.Dockerfile .. && docker run -d --network=opg-net -p ${port}:8002 --name colors-back-core colors-back-core >> ~/docker.log && echo [SERVER-CORE] Server is now running at: ${port}