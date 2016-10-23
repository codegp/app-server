#!/bin/bash

set -e

PROJECT_ID=$(gcloud config list project --format "value(core.project)" 2> /dev/null)
PUSH=false
TAG="latest"

while [ "$1" != "" ]; do
    PARAM=`echo $1 | awk -F= '{print $1}'`
    VALUE=`echo $1 | awk -F= '{print $2}'`
    case $PARAM in
        -d |--push)
            PUSH=$VALUE
            ;;
        -t | --tag)
            TAG=$VALUE
            ;;
        *)
            echo "ERROR: unknown parameter \"$PARAM\""
            exit 1
            ;;
    esac
    shift
done

cd $GOPATH/src/github.com/codegp/app-server/server
CGO_ENABLED=0 go build
cd $GOPATH/src/github.com/codegp/app-server
docker build -t gcr.io/$PROJECT_ID/app-server:$TAG .

if $PUSH; then
  gcloud docker push gcr.io/$PROJECT_ID/app-server:$TAG
fi
