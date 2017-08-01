#!/usr/bin/env bash

env

cd $PIPELINES

for yml in *
do
    name=$(echo $yml | cut -f 1 -d '.')
    echo name: $name file: $yml
done