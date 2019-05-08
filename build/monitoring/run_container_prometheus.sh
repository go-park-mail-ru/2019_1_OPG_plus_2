#!/usr/bin/env bash

docker build --rm -t colors-mon-prometheus -f ../dockerfiles/prometheus.Dockerfile ./../.. &&
docker run -d \
           --network=opg-net \
           --network-alias=colors-mon-prometheus \
           -p 9090:9090 \
           --name colors-mon-prometheus colors-mon-prometheus >> ~/docker.log