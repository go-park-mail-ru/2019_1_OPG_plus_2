#!/usr/bin/env bash

ssh-keyscan -H $PRODUCTION_MACHINE_ADDRESS >> ~/.ssh/known_hosts
chmod 600 ./deploy_key_production
ssh -i ./deploy_key_production $PRODUCTION_MACHINE_USERNAME@$PRODUCTION_MACHINE_ADDRESS << EOF
cd back
git checkout production
git fetch
cd build
./clear_up_containers.sh && ./db_run_container.sh && ./auth_run_container.sh && ./server_run_container.sh
exit
EOF
