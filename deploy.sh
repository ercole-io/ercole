#!/bin/bash

docker tag sorintlab/ercole-services sorintlab/ercole-services:${VERSION}
[[ $? == 0 ]] && echo "$(echo ${DOCKER_PASSWORD} | base64 -d)" | docker login --username "$(echo ${DOCKER_USERNAME} | base64 -d)" --password-stdin
[[ $? == 0 ]] && docker push sorintlab/ercole-services:${VERSION}

[[ $? == 0 ]] && curl --location --request POST '${REDEPLOY_URL}' --header 'Content-Type: application/json' --data-raw '{ "namespace": "ercole", "podname" : "ercole-services" }'
