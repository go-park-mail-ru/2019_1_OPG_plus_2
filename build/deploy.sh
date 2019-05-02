#!/usr/bin/env bash

ssh -i ./deploy_key $COLORS_MACHINE_USERNAME@$COLORS_MACHINE_ADDRESS

cd go-park-mail-ru/2019_1_OPG_plus_2
cd build
./db_run_container.sh && ./auth_run_container.sh && ./server_run_container.sh
exit