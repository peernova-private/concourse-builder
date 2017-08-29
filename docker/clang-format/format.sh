#!/usr/bin/env bash
BUILD_DIR=`pwd`

set -ex

CHECK_ARGS=true

if [ -z "$SOURCE"  ]
then
  echo "Please specify TARGET env variable"
  echo "It specifies the place the code needs formatting is"
  CHECK_ARGS=false
fi

if [ -z "$TARGET"  ]
then
  echo "Please specify TARGET env variable"
  echo "It specifies where the formatted code to be placed"
  CHECK_ARGS=false
fi

if [ -z "$EXTENSIONS"  ]
then
  echo "Please specify EXTENSION env variable"
  echo "It specifies the extensions of the files that need to be formatted"
  CHECK_ARGS=false
fi

if [ "$CHECK_ARGS" == "false" ]
then
    exit 1
fi

mkdir -p $TARGET

cp -R $SOURCE/. $TARGET/
cd $TARGET

find . -name $EXTENSIONS | xargs clang-format -i