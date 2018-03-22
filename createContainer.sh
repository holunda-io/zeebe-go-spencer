#!/usr/bin/env bash

APP_NAME=$1

cd $APP_NAME

echo "Building go app..."
go build

# Copy zbctl to working directory as it's used inside docker containers for zeebe health check
cp $(which zbctl) .
echo "Building docker container..."
docker build . -t $APP_NAME:latest
rm zbctl
