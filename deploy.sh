#!/bin/bash

FPC_PATH=$(pwd)
export FPC_PATH

make -C $FPC_PATH/utils/docker pull

make -C $FPC_PATH build
make -C $FPC_PATH/ercc docker

make -C $FPC_PATH/cc


make -C $FPC_PATH/application build docker
