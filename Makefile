export CGO_ENABLED?=0
DOCKER_TAG?=dev
DOCKER_IMAGE?=installer

default: build

verify:
	GO111MODULE=on go mod verify

build:
	go build -v -ldflags '-s -w' ./cmd/installer

docker:
	docker build -t "$(DOCKER_IMAGE):$(DOCKER_TAG)" .

release: build docker
	docker push "$(DOCKER_IMAGE):$(DOCKER_TAG)"
