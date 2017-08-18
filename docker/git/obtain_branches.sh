#!/usr/bin/env bash
BUILD_DIR=`pwd`

set -ex

CHECK_ARGS=true

if [ -z "$GIT_REPO_DIR"  ]
then
  echo "Please specify GIT_REPO_DIR env variable"
  echo "It specifies the relative to the current directory path of the git repo"
  CHECK_ARGS=false
fi

if [ -z "$GIT_PRIVATE_KEY"  ]
then
  echo "Please specify GIT_PRIVATE_KEY env variable"
  echo "A private key to authenticate to the repo git repo"
  CHECK_ARGS=false
fi

if [ -z "$OUTPUT_DIR"  ]
then
  echo "Please specify OUTPUT_DIR env variable"
  echo "It specifies a relative to the current directory path where the result file to be written"
  CHECK_ARGS=false
fi

if [ "$CHECK_ARGS" == "false" ]
then
    exit 1
fi


mkdir -p $BUILD_DIR/$OUTPUT_DIR && cd $BUILD_DIR/$GIT_REPO_DIR

eval `ssh-agent -s`
mkdir -p /root/.ssh
ssh-keyscan github.com >> /root/.ssh/known_hosts
echo "$GIT_PRIVATE_KEY" > $BUILD_DIR/private_key
chmod 600 $BUILD_DIR/private_key
ssh-add $BUILD_DIR/private_key

git config remote.origin.fetch "+refs/heads/*:refs/remotes/origin/*"
git fetch origin

git branch -r | sed 's/.*-\>//g' | sed 's/^ *origin\///g' | sort | uniq > $BUILD_DIR/$OUTPUT_DIR/branches

rm $BUILD_DIR/private_key

cat $BUILD_DIR/$OUTPUT_DIR/branches
