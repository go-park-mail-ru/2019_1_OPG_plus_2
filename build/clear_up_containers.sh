#!/usr/bin/env bash


clear_cache="false"

while getopts 'c' flag; do
  case "${flag}" in
    c) clear_cache="true" ;;
  esac
done

docker stop $(cat ~/docker.log)
docker rm $(cat ~/docker.log)
if [[ ${clear_cache} = "true" ]]
 then
 echo "REMOVING CACHE"
docker image prune -f
fi

rm ~/docker.log
echo
echo "Images left (just for checking whether no trash containers and images left)"
docker images
echo
docker ps -a