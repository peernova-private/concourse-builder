#!/usr/bin/env bash
BUILD_DIR=`pwd`

set -ex

CHECK_ARGS=true

if [ -z "$USER_NAME" ]
then
  echo "Please specify USER_NAME env variable"
  echo "It specifies the user name for the commit"
  CHECK_ARGS=false
fi

if [ -z "$USER_EMAIL" ]
then
  echo "Please specify USER_EMAIL env variable"
  echo "It specifies the user email for the commit"
  CHECK_ARGS=false
fi

if [ -z "$BRANCH" ]
then
  echo "Please specify BRANCH env variable"
  echo "It specifies the branch in which the commit to be made"
  CHECK_ARGS=false
fi

if [ -z "$MESSAGE" ]
then
  echo "Please specify MESSAGE env variable"
  echo "It specifies the message to use for the commit"
  CHECK_ARGS=false
fi


if [ "$CHECK_ARGS" == "false" ]
then
    exit 1
fi

cp -R $SOURCE/. $TARGET/
cd $TARGET

git add .

set +e
if ! git diff --cached --exit-code
then
    set -e
    git checkout -b "$BRANCH"

    git config --global user.email "$USER_EMAIL"
    git config --global user.name "$USER_NAME"

    git commit --message="$MESSAGE"
fi