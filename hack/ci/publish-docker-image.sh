#!/usr/bin/env bash

set -euo pipefail
cd $(dirname $0)/../..

echodate() {
  echo "[$(date -Is)]" "$@"
}

repo=quay.io/kubermatic/kubermatic-installer

image="$repo:$(git rev-parse HEAD)"

echo "Building Docker image $image"
docker build -t "$image" .
docker push -t "$image"

tag="$(git tag -l --points-at HEAD)"
if [ -n "$tag" ]; then
  taggedImage="$repo:$tag"
  echo "Re-tagging $image as $taggedImage"
  docker tag "$image" "$taggedImage"
  docker push "$taggedImage"
fi

echodate "Build finished."
