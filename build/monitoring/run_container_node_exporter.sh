#!/usr/bin/env bash

docker run -d \
           -p 9100:9100 \
           --network=opg-net \
           --network-alias=colors-node-exporter \
           --name colors-node-exporter prom/node-exporter >> ~/docker.log