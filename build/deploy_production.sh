#!/usr/bin/env bash

ssh-keyscan -H $PRODUCTION_MACHINE_ADDRESS >> ~/.ssh/known_hosts
chmod 600 ./deploy_key_production
ssh -i ./deploy_key_production $PRODUCTION_MACHINE_USERNAME@$PRODUCTION_MACHINE_ADDRESS << EOF
cd back
git fetch
git checkout production
git fetch
cd build

./initial.sh
./clear_up_containers.sh -c

docker run -d --name "opg-net-proxyhost" \
  --cap-add=NET_ADMIN --cap-add=NET_RAW \
  --restart on-failure \
  --net=opg-net --network-alias 'proxyhost' \
  qoomon/docker-host >> ~/docker.log

./auth_run_container.sh -m production && ./server_run_container.sh -m production
exit
EOF
