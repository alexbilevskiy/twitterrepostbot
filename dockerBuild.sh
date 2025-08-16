#!/bin/bash
if [[ $DOCKER_REGISTRY == "" ]];
  then
    echo "specifiy DOCKER_REGISTRY env"
    exit
  else
    echo "using registry ${DOCKER_REGISTRY}"
fi

PROJECT=twitterrepostbot

set -x
go mod vendor
docker build -t ${DOCKER_REGISTRY}/${PROJECT}:latest . && docker push ${DOCKER_REGISTRY}/${PROJECT}:latest