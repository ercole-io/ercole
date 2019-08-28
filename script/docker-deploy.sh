#!/bin/sh
echo "$QUAY_BOT_PASSWORD" | docker login -u "$QUAY_BOT_USERNAME" --password-stdin quay.io
docker push quay.io/amreo/ercole-server