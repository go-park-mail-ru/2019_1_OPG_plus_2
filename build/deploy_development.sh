#!/usr/bin/env bash

ssh-keyscan -H $DEVELOPMENT_MACHINE_ADDRESS >> ~/.ssh/known_hosts
chmod 600 ./deploy_key_development
ssh -i ./deploy_key_development $DEVELOPMENT_MACHINE_USERNAME@$DEVELOPMENT_MACHINE_ADDRESS << EOF
cd go-park-mail-ru/2019_1_OPG_plus_2
git fetch
git checkout dev
git fetch
cd build
./initial.sh
./clear_up_containers.sh
./db_run_container.sh && ./auth_run_container.sh && ./server_run_container.sh
exit
EOF
