#!/usr/bin/env bash
BUILD_DIR=`pwd`

set -ex

. /bin/fly/authenticate.sh

if [ -z "$PIPELINE_REGEX"  ]
then
  echo "Please specify PIPELINE_REGEX env variable"
  echo "It specifies a regexp pattern which pipelines to be considered for removal"
  exit 1
fi

cd $PIPELINES

CURRENT_PIPELINES=$(fly --target trgt pipelines | awk '{ print $1 }' | sort)
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

    fly -t trgt destroy-pipeline --non-interactive --pipeline=$pipeline
done
