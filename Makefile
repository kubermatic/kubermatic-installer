default: build

.PHONY: genassets
genassets:
	cd install-wizard && make prod
	go run pkg/assets/generate/generate.go

.PHONY: wizard
wizard:
	go run -tags=dev cmd/installer/main.go wizard

.PHONY: build
build: genassets
	go build -v -ldflags '-s -w' ./cmd/installer

.PHONY: docker
docker:
	CGO_ENABLED=0 make build
	docker build -t installer .
