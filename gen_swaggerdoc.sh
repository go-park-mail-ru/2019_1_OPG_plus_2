#!/bin/bash

dir=("cmd/core/main.go")

for var in ${dir[@]}
do
    swag init -g ${var}
done
