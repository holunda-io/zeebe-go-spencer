#!/bin/bash

cd $GOPATH
mkdir -p src/github.com/zeebe-io/
cd src/github.com/zeebe-io/

echo "Checkout zbc-go in version 0.7.0"
git clone git@github.com:zeebe-io/zbc-go.git
cd zbc-go
git checkout 0.7.0

echo "Build and install zbc-go global"
make build
sudo make install