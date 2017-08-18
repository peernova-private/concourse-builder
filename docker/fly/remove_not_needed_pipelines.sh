#!/usr/bin/env bash

set -ex


if [ ! -z "$INSECURE" ]
then
  INSECUR_VAR=" -k"
fi

fly --target trgt login$INSECUR_VAR --concourse-url $CONCOURSE_URL --username $CONCOURSE_USER --password $CONCOURSE_PASSWORD

if [ -z "$PIPELINE_REGEX"  ]
then
  echo "Please specify PIPELINE_REGEX env variable"
  echo "It specifies the regexp pattern to wortk with pipelens"
  exit 1
fi

cd $PIPELINES

EXIST_PIPELINES=$(fly --target trgt$INSECUR_VAR pipelines | awk '{ print $1 }' | sort)
PIPELINE_FILES=$(for yml in *; do name=$(echo $yml | cut -f 1 -d '.');done)

for pipeline in $EXIST_PIPELINES
do
    if [[ $pipeline =~ $PIPELINE_REGEX ]]; then
        if ! echo "$PIPELINE_FILES" | grep "$pipeline"; then
            fly -t trgt$INSECUR_VAR destroy-pipeline --pipeline=$pipeline --non-interactive
        fi
    fi
done
