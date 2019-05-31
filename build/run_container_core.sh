#!/usr/bin/env bash

port=8002
colors_mode="IN_DOCKER_NET"
db_mode="IN_DOCKER_NET"

while getopts 'm:p:d:' flag; do
  case "${flag}" in
    p) port=${OPTARG};;
    m) mode=${OPTARG}
        case ${mode} in
            dockernet) colors_mode="IN_DOCKER_NET";;
            production) colors_mode="PRODUCTION";;
        esac
        ;;
     d) mode=${OPTARG}
        case ${mode} in
            dockernet) db_mode="IN_DOCKER_NET";;
            dockerdb) db_mode="USE_DOCKER_DB";;
            production) db_mode="PRODUCTION";;
        esac
        ;;
  esac
done

echo [SERVER-CORE] The published port of the container is ${port}

docker build --rm -t colors-back-core -f dockerfiles/core.Dockerfile .. && \
docker run -d \
           -e COLORS_SERVICE_USE_MODE=${colors_mode} \
           -e COLORS_DB=${db_mode} \
           -e DB_HOST=${DB_HOST} \
           -e DB_PORT=${DB_PORT} \
           -e DB_USERNAME=${DB_USERNAME} \
           -e DB_PASSWORD=${DB_PASSWORD} \
           --network=opg-net \
           --network-alias=colors-back-core \
           --mount type=bind,source=$PWD/../upload,target=/root/colors-core-service/upload \
           -p ${port}:8002 \
           --name colors-back-core colors-back-core >> ~/docker.log && \
echo [SERVER-CORE] Server is now running at: ${port}
exit 0
