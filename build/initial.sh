#!/usr/bin/env bash

if [[ ! $(docker network ls | grep opg-net) ]]
then
echo Creating required docker network opg-net...
docker network create opg-net
else
echo Docker network requirements are met
fi

