#!/usr/bin/env bash

if [[ "$#" -eq 1 ]]
then
    port=$1
    else
    port=50242
    fi

echo [SERVER-AUTH] The published port of the container is ${port}

docker build --rm -t colors-back-auth -f auth.Dockerfile .. && docker run -d -e COLORS_AUTH_MODE=$1--network=opg-net -p ${port}:50242 --name colors-back-auth colors-back-auth >> ~/docker.log && echo [SERVER-AUTH] Server is now running at: ${port}