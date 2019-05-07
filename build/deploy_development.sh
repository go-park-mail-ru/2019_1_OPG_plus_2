#!/usr/bin/env bash

ssh-keyscan -H $DEVELOPMENT_MACHINE_ADDRESS >> ~/.ssh/known_hosts
chmod 600 ./deploy_key_development
ssh -i ./deploy_key_development $DEVELOPMENT_MACHINE_USERNAME@$DEVELOPMENT_MACHINE_ADDRESS << EOF
if [ -f ~/.bash_exports ]; then
    . ~/.bash_exports
fi

cd go-park-mail-ru/2019_1_OPG_plus_2
git checkout dev
git pull

cd build
./initial.sh
./clear_up_containers.sh -c
./run_container_db.sh && ./run_container_auth.sh && ./server_run_container.sh
exit
EOF
