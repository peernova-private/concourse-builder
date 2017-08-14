#!/usr/bin/env bash

set -ex

fly --target trgt login --insecure --concourse-url $CONCOURSE_URL --username $CONCOURSE_USER --password $CONCOURSE_PASSWORD

cd $PIPELINES

for yml in *
do
    name=$(echo $yml | cut -f 1 -d '.')
    # TODO: check if the pipeline if the pipeline exist.
    # If exists we are good, we should not update it, because it might be already in better
    # shape, based on self-update feature.
    # If it does not exists, we need to create it.
done