#!/usr/bin/env bash

set -ex

INSECURE="1"

if [ ! -z "$INSECURE" ]
then
  INSECUR_VAR=" -k"
fi

FLY_VERSION=`fly -v`
CONCOURSE_VERSION=`curl $CONCOURSE_URL/api/v1/info$INSECUR_VAR | awk -F ',' ' { print $1 } ' | awk -F ':' ' { print $2 } ' | sed -e 's/^"//' -e 's/"$//'`

if [ ! "$FLY_VERSION" == "$CONCOURSE_VERSION" ]
then
    exit 1
fi
