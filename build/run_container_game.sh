#!/usr/bin/env bash

port=8004
colors_mode="IN_DOCKER_NET"
db_mode="IN_DOCKER_NET"

while getopts 'm:tp:d:' flag; do
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


echo [MODE] ${colors_mode}
echo [SERVER-GAME] The published port of the container is ${port}

docker build --rm -t colors-back-game -f dockerfiles/game.Dockerfile .. && \
docker run -d \
           -e COLORS_SERVICE_USE_MODE=${colors_mode} \
           -e COLORS_DB=${db_mode} \
           -e DB_HOST=${DB_HOST} \
           -e DB_PORT=${DB_PORT} \
           -e DB_USERNAME=${DB_USERNAME} \
           -e DB_PASSWORD=${DB_PASSWORD} \
           --network=opg-net \
           --network-alias=colors-back-game \
           --restart=always \
           -p ${port}:8004 \
           --name colors-back-game colors-back-game >> ~/docker.log && \
echo [SERVER-GAME] Server is now running at: ${port}
