#!/bin/bash

# basic stress test

set -e

while true; do
	go get -u -v github.com/wallyworld/core/utils
	export GOMAXPROCS=$[ 1 + $[ RANDOM % 128 ]]
        go test github.com/wallyworld/core/... 2>&1
done
