#!/usr/bin/env bash

set -ex

if [ -z "$CONCOURSE_URL"  ]
then
  echo "Please specify CONCOURSE_URL env variable"
  echo "It specifies the url to the concourse"
  exit 1
fi

if [ ! -z "$INSECURE" ]
then
  INSECURE_VAR=" --insecure"
fi

FLY_VERSION=`fly --version`
CONCOURSE_VERSION=`curl $CONCOURSE_URL/api/v1/info$INSECURE_VAR |\
    awk -F ',' ' { print $1 } ' |\
    awk -F ':' ' { print $2 } ' |\
    sed -e 's/^"//' -e 's/"$//'`

if [ ! "$FLY_VERSION" == "$CONCOURSE_VERSION" ]
then
    exit 1
fi
