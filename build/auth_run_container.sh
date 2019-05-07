#!/usr/bin/env bash

port=50242
auth_mode="RUNTIME"
colors_mode="IN_DOCKER_NET"
db_mode="IN_DOCKER_NET"

while getopts 'm:tp:d:' flag; do
  case "${flag}" in
    t) auth_mode="TEST"
        echo [TESTING];;
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


echo [AUTH] ${auth_mode}
echo [MODE] ${colors_mode}
echo [SERVER-AUTH] The published port of the container is ${port}

docker build --rm -t colors-back-auth -f auth.Dockerfile .. && \
docker run -d \
    -e COLORS_SERVICE_USE_MODE=${colors_mode} \
    -e COLORS_AUTH_MODE=${auth_mode} \
    -e COLORS_DB=${db_mode} \
    -e DB_HOST=${DB_HOST} \
    -e DB_PORT=${DB_PORT} \
    -e DB_USERNAME=${DB_USERNAME} \
    -e DB_PASSWORD=${DB_PASSWORD} \
    --network=opg-net \
    -p ${port}:50242 \
    --name colors-back-auth colors-back-auth \>> ~/docker.log && \
echo [SERVER-AUTH] Server is now running at: ${port}
