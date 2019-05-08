#!/usr/bin/env bash

ssh-keyscan -H $PRODUCTION_MACHINE_ADDRESS >> ~/.ssh/known_hosts
chmod 600 ./deploy_key_production
ssh -i ./deploy_key_production $PRODUCTION_MACHINE_USERNAME@$PRODUCTION_MACHINE_ADDRESS << EOF
if [ -f ~/.bash_exports ]; then
    . ~/.bash_exports
fi

cd back
git checkout production
git pull

cd build
./initial.sh
./clear_up_containers.sh -c
./run_container_cookie.sh && ./run_container_auth.sh -d production && ./run_container_server.sh -d production
cd ../monitoring
./run_container_node_exporter.sh -d & ./run_container_prometheus.sh -d
exit
EOF
