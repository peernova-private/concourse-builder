#!/usr/bin/env bash

set -e

CHECK_ARGS=true

if [ -z "$CONCOURSE_URL"  ]
then
  echo "Please specify CONCOURSE_URL env variable"
  echo "It specifies the url to the concourse"
  CHECK_ARGS=false
fi

if [ -z "$CONCOURSE_TEAM"  ]
then
  echo "Please specify CONCOURSE_TEAM env variable"
  echo "It specifies the team to the concourse"
  CHECK_ARGS=false
fi

if [ -z "$CONCOURSE_USER" -o -z "$CONCOURSE_PASSWORD"  ]
then
  if [ -z "$CONCOURSE_TOKEN" ]
  then
    echo "Please specify (CONCOURSE_USER and CONCOURSE_PASSWORD) or CONCOURSE_TOKEN env variable"
    echo "It specifies a user to authenticate in concourse"
    CHECK_ARGS=false
  else
    AUTHENTICATE="--token $CONCOURSE_TOKEN"
  fi
else
  AUTHENTICATE="--username $CONCOURSE_USER --password $CONCOURSE_PASSWORD"
fi

if [ "$CHECK_ARGS" == "false" ]
then
    exit 1
fi

if [ ! -z "$INSECURE" ]
then
  INSECURE_VAR=" --insecure"
fi

fly --target trgt login$INSECURE_VAR \
    --concourse-url $CONCOURSE_URL \
    --team-name $CONCOURSE_TEAM \
    $AUTHENTICATE
