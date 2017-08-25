#!/usr/bin/env bash

set -ex

. /bin/fly/authenticate.sh

cd $PIPELINES

CURRENT_PIPELINES=$(fly --target trgt pipelines -a | awk '{ print $1 }')

for yml in *
do
    name=$(echo $yml | cut -f 1 -d '.')

    if echo $CURRENT_PIPELINES | grep -w $name > /dev/null
    then
        echo "'$name' pipeline already exists, skipping"
    else
        fly --target trgt set-pipeline --non-interactive --pipeline=$name --config=$yml && echo "$name created"
    fi
done

