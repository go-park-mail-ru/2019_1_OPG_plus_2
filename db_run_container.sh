#!/usr/bin/env bash
if [[ "$#" -eq 1 ]]
then
    port=$1
    else
    port=12345
    fi

echo [DATABASE] The published port of the container is ${port}

docker build --rm -t colors-database -f db.Dockerfile .
docker run -d --network=opg-net -p ${port}:3306 --rm --name colors-db colors-database >> docker.log

echo [DATABASE] Database is now running at: ${port}