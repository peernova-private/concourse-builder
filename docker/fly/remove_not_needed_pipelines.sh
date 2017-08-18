#!/usr/bin/env bash
BUILD_DIR=`pwd`

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

if [ -z "$PIPELINE_REGEX"  ]
then
  echo "Please specify PIPELINE_REGEX env variable"
  echo "It specifies a regexp pattern which pipelines to be considered for removal"
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

CURRENT_PIPELINES=$(fly --target trgt$INSECUR_VAR pipelines | awk '{ print $1 }' | sort)
PIPELINE_FILES=$(for yml in *; do echo $(echo $yml | cut -f 1 -d '.'); done)

for pipeline in $CURRENT_PIPELINES
do
    # check if the pipeline matches the regular expression
    if [[ ! $pipeline =~ $PIPELINE_REGEX ]]
    then
        continue
    fi

    # check if the pipeline is in the list of generated pipelines
    if echo "$PIPELINE_FILES" | grep "$pipeline"
    then
        continue
    fi

    fly -t trgt$INSECUR_VAR destroy-pipeline --non-interactive --pipeline=$pipeline
done
