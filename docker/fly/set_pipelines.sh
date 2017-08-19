#!/usr/bin/env bash

set -ex

CHECK_ARGS=true

if [ -z "$CONCOURSE_URL"  ]
then
  echo "Please specify CONCOURSE_URL env variable"
  echo "It specifies the url to the concourse"
  CHECK_ARGS=false
fi

if [ -z "$CONCOURSE_USER"  ]
then
  echo "Please specify CONCOURSE_USER env variable"
  echo "It specifies a user to authenticate in concourse"
  CHECK_ARGS=false
fi

if [ -z "$CONCOURSE_PASSWORD"  ]
then
  echo "Please specify CONCOURSE_PASSWORD env variable"
  echo "It specifies the user password to use to authenticate in concourse"
  CHECK_ARGS=false
fi

if [ "$CHECK_ARGS" == "false" ]
then
    exit 1
fi

if [ ! -z "$INSECURE" ]
then
  INSECURE_VAR=" --insecure"
fi

fly --target trgt login$INSECURE_VAR --concourse-url $CONCOURSE_URL --username $CONCOURSE_USER --password $CONCOURSE_PASSWORD

cd $PIPELINES

for yml in *
do
    name=$(echo $yml | cut -f 1 -d '.')
    fly -t trgt set-pipeline --non-interactive --pipeline=$name --config=$yml
done