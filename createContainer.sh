#!/usr/bin/env bash

APP_NAME=$1
ORIGINAL_PWD=$(pwd)

cd $APP_NAME

echo "Building go app..."
go build

cd $ORIGINAL_PWD

# Copy zbctl to working directory as it's used inside docker containers for zeebe health check
cp $(which zbctl) .
echo "Building docker container..."
docker build . -f $APP_NAME/Dockerfile -t $APP_NAME:latest
rm zbctl
