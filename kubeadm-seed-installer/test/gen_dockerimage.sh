#!/usr/bin/env bash

NAME=quay.io/kubermatic/seed-installer-e2e-builder
TAG=0.5

docker build -t $NAME:$TAG .
docker push $NAME:$TAG
