#!/usr/bin/env bash

docker="false"
while getopts 'd' flag; do
  case ${flag} in
    d) docker="true";;
  esac
done

docker build --rm -t colors-mon-prometheus -f ../dockerfiles/prometheus.Dockerfile ./../.. &&

if [[ docker = "true" ]]
then
docker run -d \
           -p 9100:9100 \
           --network=opg-net \
           --network-alias=colors-node-exporter \
           --name colors-node-exporter prom/node-exporter >> ~/docker.log
else
docker run --rm \
           --network=host \
           --name colors-node-exporter prom/node-exporter
fi
