#!/usr/bin/env bash

set -e

NAME=quay.io/kubermatic/seed-installer-e2e-builder
#TAG=0.6
TAG=pkwip

docker build -t $NAME:$TAG ../ -f ./Dockerfile
docker push $NAME:$TAG
