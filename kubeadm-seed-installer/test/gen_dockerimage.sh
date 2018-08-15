#!/usr/bin/env bash
# This file should only used for debugging purposes since the ci build process happens in .drone.yml

set -e

NAME=quay.io/kubermatic/seed-installer-e2e-builder
TAG=wip

docker build -t $NAME:$TAG ../ -f ./Dockerfile
docker push $NAME:$TAG
