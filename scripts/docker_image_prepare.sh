#!/usr/bin/env bash

set -ex

mkdir -p prepared
cp $DOCKERFILE_DIR/* prepared

cd  prepared

echo FROM $FROM_IMAGE > Dockerfile
echo >> Dockerfile

if [ ! -z "$EVAL" ]
then
    eval "$EVAL" >> Dockerfile
fi

echo >> Dockerfile
cat ../$DOCKERFILE_DIR/steps >> Dockerfile

cat Dockerfile