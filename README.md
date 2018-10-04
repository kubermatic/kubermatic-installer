# Kubermatic-installer

This repository contains Helm charts for installing Kubermatic and the
seed installer which can be used to set up a master or seed cluster.

Please refer to [our docs](https://docs.kubermatic.io) for more information.

## Hacking

The installer consists of two parts: A backend in Go and a frontend
written in Angular 6. During development you will want to use Angular's
internal webserver to deliver the assets, while in production the
assets will be compiled and bundled into the Go binary.

### Frontend

The frontend lives in `install-wizard/` and has a Dockerfile that
contains the required Node.js toolset. You can use the `Makefile`
to build a local Docker image via `make docker` and then you
most likely want to start the container with a shell and use the
various dev commands:

* `make shell` starts the container and drops you into a shell.
  * `ng serve` starts Angular's dev server, afterwards you can access
    the wizard at http://localhost:4200.
* `make lint` will run tslint on the wizard.

### Backend

Dependencies are already vendored, to make changes use Go 1.11's
module commands and commit new libraries to the repository.

The backend binary lives in `cmd/installer/`, so use whatever Go
build command you like, like `go build -v ./cmd/installer`.
Afterwards you can run `./installer help` and see what it offers.

There is a `make build` command, but that will also compile the
Angular assets and is too slow for repeated builds while hacking.
