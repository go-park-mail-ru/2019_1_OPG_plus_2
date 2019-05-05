#!/usr/bin/env bash

port=50242
auth_mode="RUNTIME"
colors_mode="IN_DOCKER_NET"

while getopts 'm:tp:' flag; do
  case "${flag}" in
    t) auth_mode="TEST"
        echo [TESTING];;
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

echo [MODE] ${colors_mode}
echo [SERVER-AUTH] The published port of the container is ${port}

docker build --rm -t colors-back-auth -f auth.Dockerfile .. &&\
docker run -d -e COLORS_AUTH_MODE=${auth_mode} -e COLORS_SERVICE_USE_MODE=${colors_mode} --network=opg-net\
 -p ${port}:50242 --name colors-back-auth colors-back-auth >> ~/docker.log\
&& echo [SERVER-AUTH] Server is now running at: ${port}