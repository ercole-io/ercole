#!/bin/bash

docker tag ercole/ercole-services ercole/ercole-services:${VERSION}
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push ercole/ercole-services:${VERSION}