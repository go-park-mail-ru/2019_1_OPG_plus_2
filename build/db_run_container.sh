#!/usr/bin/env bash

port=12345
while getopts 'p:' flag; do
  case ${flag} in
    p) port=${OPTARG};;
  esac
done


echo [DATABASE] The published port of the container is ${port}

docker build --rm -t colors-db -f db.Dockerfile .. && \
docker run -d --network=opg-net -p ${port}:3306 --name colors-db colors-db >> ~/docker.log && \
echo [DATABASE] Database is now running at: ${port}
