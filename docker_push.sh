#!/bin/bash

docker tag ercole/ercole-services ercole/ercole-services:${VERSION}
[[ $? == 0 ]] && echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
[[ $? == 0 ]] && docker push ercole/ercole-services:${VERSION}