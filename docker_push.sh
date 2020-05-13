#!/bin/bash

docker tag ercole/ercole-services ercole/ercole-services:${VERSION}
[[ $? == 0 ]] && echo "$(echo ${DOCKER_PASSWORD} | base64 -d)" | docker login --username "$(echo ${DOCKER_USERNAME} | base64 -d)" --password-stdin
[[ $? == 0 ]] && docker push ercoleorg/ercole-services:${VERSION}