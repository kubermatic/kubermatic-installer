default: build

.PHONY: genassets
genassets:
	go run pkg/assets/generate/generate.go

.PHONY: wizard
wizard:
	go run -tags=dev cmd/installer/main.go wizard

.PHONY: build
build: genassets
	go build -v ./cmd/installer
