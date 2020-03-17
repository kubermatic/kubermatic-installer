#!/usr/bin/env bash

set -euo pipefail
cd $(dirname $0)/../..

source ./hack/lib.sh

if [ -n "${QUAY_IO_USERNAME:-}" ]; then
  echodate "Logging into Quay"
  docker ps > /dev/null 2>&1 || start-docker.sh
  retry 5 docker login -u "$QUAY_IO_USERNAME" -p "$QUAY_IO_PASSWORD" quay.io
  echodate "Successfully logged into Quay"
fi

repo=quay.io/kubermatic/installer

image="$repo:$(git rev-parse HEAD)"

echodate "Building Docker image $image"
docker build -t "$image" .
docker push "$image"

tag="$(git tag -l --points-at HEAD)"
if [ -n "$tag" ]; then
  taggedImage="$repo:$tag"
  echodate "Re-tagging $image as $taggedImage"
  docker tag "$image" "$taggedImage"
  docker push "$taggedImage"
fi

echodate "Build finished."
