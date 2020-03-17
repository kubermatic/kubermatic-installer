export GO111MODULE=on
GIT_HASH?=$(shell git rev-parse HEAD)
INSTALLER_VERSION?=$(shell git tag -l --points-at HEAD)

ifeq "$(INSTALLER_VERSION)" ""
	INSTALLER_VERSION=$(GIT_HASH)
endif

default:
	go build -v -mod vendor -ldflags '-s -w -X github.com/kubermatic/kubermatic-installer/pkg/shared.INSTALLER_VERSION=$(INSTALLER_VERSION) -X github.com/kubermatic/kubermatic-installer/pkg/shared.INSTALLER_GIT_HASH=$(GIT_HASH)' ./cmd/installer

verify:
	go mod verify

vendor:
	go mod vendor
