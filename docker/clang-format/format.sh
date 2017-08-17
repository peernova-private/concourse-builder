#!/usr/bin/env bash

set -ex

if [ -z "$EXTENSION"  ]
then
  echo "Please specify EXTENSION env variable"
  exit 1
fi

find . -name "*.$EXTENSION" | xargs clang-format -i