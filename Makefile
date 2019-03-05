export CGO_ENABLED?=0
DOCKER_TAG?=dev
DOCKER_IMAGE?=installer
HAS_NPM:=$(shell command -v npm 2> /dev/null)

default: build

assets:
ifdef HAS_NPM
	cd install-wizard && make build
else
	cd install-wizard && CMD="make build" make shell
endif
	go run pkg/assets/generate/generate.go

wizard:
	go run -tags=dev cmd/installer/main.go wizard

verify:
	GO111MODULE=on go mod verify

build:
	go build -v -ldflags '-s -w' ./cmd/installer

docker:
	docker build -t "$(DOCKER_IMAGE):$(DOCKER_TAG)" .

release: assets build docker
	docker push "$(DOCKER_IMAGE):$(DOCKER_TAG)"
