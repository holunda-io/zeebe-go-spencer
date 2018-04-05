#!/usr/bin/env bash

echo "Waiting for Zeebe to be available..."
for ((i=1; $i<5; i++)); do
    if $(./zbctl --config config.toml describe topology > /dev/null 2>&1); then
        echo "Zeebe is Ready!"
        break;
    fi
    echo "Still waiting..."
    sleep 1
done

APPNAME=$1

echo "Zeebe is ready (or wait time expired), starting $APPNAME..."
./$APPNAME
