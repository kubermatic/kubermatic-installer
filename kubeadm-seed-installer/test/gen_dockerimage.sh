#!/usr/bin/env bash

set -e

NAME=quay.io/kubermatic/seed-installer-e2e-builder
TAG=0.6

docker build -t $NAME:$TAG .
docker push $NAME:$TAG
