#!/usr/bin/env bash

set +ex

mkdir -p prepared
cd  prepared

echo FROM $FROM_IMAGE > Dockerfile
echo >> Dockerfile
cat ../$DOCKER_STEPS >> Dockerfile
