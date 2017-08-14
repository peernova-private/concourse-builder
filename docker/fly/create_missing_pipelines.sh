#!/usr/bin/env bash

set -ex

fly --target trgt login --insecure --concourse-url $CONCOURSE_URL --username $CONCOURSE_USER --password $CONCOURSE_PASSWORD

cd $PIPELINES

EXIST_PIPELINES=$(fly --target trgt pipelines -a | awk '{ print $1 }')

for yml in *
do
    name=$(echo $yml | cut -f 1 -d '.')

    if echo $EXIST_PIPELINES | grep -w $name
    then
        echo "'$name' pipeline already exists, skipping"
    else
        fly --target trgt --pipeline $name -c $yml --non-interactive && echo "$name created"
    fi
done