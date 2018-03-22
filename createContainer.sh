#!/usr/bin/env bash

APP_NAME=$1

cd $APP_NAME

echo "Building go app..."
go build

echo "Building docker container..."
docker build . -t $APP_NAME:latest
