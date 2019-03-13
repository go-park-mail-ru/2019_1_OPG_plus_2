#!/usr/bin/env bash

if test ! -z "$(docker images -q mariadb:latest)"; then
  docker pull mariadb
fi

docker run --name mariadb -e