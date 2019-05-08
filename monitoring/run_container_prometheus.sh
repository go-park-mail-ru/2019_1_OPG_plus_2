#!/usr/bin/env bash

docker="false"
while getopts 'd' flag; do
  case ${flag} in
    d) docker="true";;
  esac
done

docker build --rm -t colors-mon-prometheus -f prometheus.Dockerfile . &&

if [[ docker = "true" ]]
then
docker run -d \
           --network=opg-net \
           --network-alias=colors-mon-prometheus \
           -p 9090:9090 \
           --name colors-mon-prometheus colors-mon-prometheus >> ~/docker.log
else
docker run --network=host \
           --rm \
           --name colors-mon-prometheus colors-mon-prometheus
fi