#!/usr/bin/env bash

port=8002
colors_mode="IN_DOCKER_NET"

while getopts 'm:p:' flag; do
  case "${flag}" in
    p) port=${OPTARG};;
    m) mode=${OPTARG}
        case ${mode} in
            dockernet) colors_mode="IN_DOCKER_NET";;
            dockerdb) colors_mode="USE_DOCKER_DB";;
            production) colors_mode="PRODUCTION";;
        esac
        ;;
  esac
done

echo [SERVER-CORE] The published port of the container is ${port}

docker build --rm -t colors-back-core -f core.Dockerfile .. &&\
docker run -d --network=opg-net -p ${port}:8002 -e COLORS_SERVICE_USE_MODE=${colors_mode}\
    --name colors-back-core colors-back-core >> ~/docker.log &&\
echo [SERVER-CORE] Server is now running at: ${port}