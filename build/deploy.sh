#!/usr/bin/env bash

ssh-keyscan -H $COLORS_MACHINE_ADDRESS >> ~/.ssh/known_hosts
chmod 600 ./deploy_key
ssh -i ./deploy_key $COLORS_MACHINE_USERNAME@$COLORS_MACHINE_ADDRESS << EOF

cd go-park-mail-ru/2019_1_OPG_plus_2
git checkout dev
git fetch
cd build
./clear_up_containers.sh && ./db_run_container.sh && ./auth_run_container.sh && ./server_run_container.sh
exit
EOF