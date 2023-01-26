#!/bin/bash
set -ex
cd `dirname $0`
docker buildx build --platform linux/amd64 -t test-transaction-service:0.0.4 --load -f ../Dockerfile-amd --target=app ../