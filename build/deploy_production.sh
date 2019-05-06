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
./auth_run_container.sh -m production && ./server_run_container.sh -m production
exit
EOF
