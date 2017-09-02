#!/usr/bin/env bash
ROOT=`pwd`

set -ex

CHECK_ARGS=true

if [ -z "$DOCKERFILE_DIR"  ]
then
  echo "Please specify DOCKERFILE_DIR env variable"
  echo "It specifies the directory where the dockerfile steps are"
  CHECK_ARGS=false
fi

if [ -z "$FROM_IMAGE"  ]
then
  echo "Please specify FROM_IMAGE env variable"
  echo "It specifies the repository to be used in the FROM clause"
  CHECK_ARGS=false
fi

if [ "$CHECK_ARGS" == "false" ]
then
    exit 1
fi


mkdir -p prepared
cp $DOCKERFILE_DIR/* prepared

for SOURCE_DIR in $SOURCE_DIRS
do
    mkdir -p prepared/$SOURCE_DIR
    cp -R $SOURCE_DIR/. prepared/$SOURCE_DIR
done

cd  prepared

echo FROM $FROM_IMAGE > Dockerfile
echo >> Dockerfile

if [ ! -z "$EVAL" ]
then
    eval "$EVAL" >> Dockerfile
fi

echo >> Dockerfile
cat $ROOT/$DOCKERFILE_DIR/steps >> Dockerfile
