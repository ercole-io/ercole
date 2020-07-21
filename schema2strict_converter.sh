#!/bin/sh
cat $1 | sed 's/"type": "object"/"type": "object",\n"additionalProperties": false/' | jq