#!/bin/bash

dir=("cmd/server/main.go")

for var in ${dir[@]}
do
    swag init -g ${var}
done
