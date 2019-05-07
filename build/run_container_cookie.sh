#!/usr/bin/env bash

port=50243
while getopts 'p:' flag; do
  case ${flag} in
    p) port=${OPTARG};;
  esac
done

echo [SERVER-COOKIE] The published port of the container is ${port}

docker build --rm -t colors-back-cookie -f dockerfiles/cookie.Dockerfile .. && \
docker run -d \
           --network=opg-net \
           -p ${port}:50243 \
           --name colors-back-cookie \
           colors-back-cookie >> ~/docker.log && \
echo [SERVER-COOKIE] Cookie-validation service is now running at: ${port}
