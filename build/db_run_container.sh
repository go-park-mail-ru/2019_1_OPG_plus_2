#!/usr/bin/env bash

if [[ "$#" -eq 1 ]]
then
    port=$1
    else
    port=12345
    fi

echo [DATABASE] The published port of the container is ${port}

docker build --rm -t colors-db -f db.Dockerfile .
docker run -d --network=opg-net -p ${port}:3306 --name colors-db colors-db >> ~/docker.log

echo [DATABASE] Database is now running at: ${port}