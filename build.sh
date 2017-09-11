#!/usr/bin/env bash
set -eu
OLD_DIR=$PWD
CURR_DIR=$(dirname $0)
cd $CURR_DIR
CGO_ENABLED=0 go build -o ./gnag
cd $OLD_DIR
