export CGO_ENABLED = 0
export GOOS = $(shell go env GOOS)
export GOARCH = $(shell go env GOARCH)

GO := go
ARCH := $(GOARCH)
ifeq ($(ARCH), arm)
	ARCH = armhf
endif

DOCKER_IMAGE := jjbubudi/tides
TAG ?= $(shell git describe --tags --exact-match 2>/dev/null)
COMMIT ?= $(shell git rev-parse --short HEAD)
VERSION = $(COMMIT)
ifneq ($(TAG),)
	VERSION = $(TAG)
endif

.PHONY: build
build:
	@$(GO) build -o dist/tides cmd/main.go

.PHONY: build-docker
build-docker: build
	@docker build --build-arg ARCH=$(ARCH) -t $(DOCKER_IMAGE):$(VERSION)-$(GOARCH) .
	@docker tag $(DOCKER_IMAGE):$(VERSION)-$(GOARCH) $(DOCKER_IMAGE):$(VERSION)
	@docker tag $(DOCKER_IMAGE):$(VERSION)-$(GOARCH) $(DOCKER_IMAGE):latest

ifeq ($(PUSH), true)
	@docker push $(DOCKER_IMAGE):$(VERSION)-$(GOARCH)
endif

.PHONY: push-manifest
push-manifest:
	@curl https://github.com/estesp/manifest-tool/releases/download/v0.9.0/manifest-tool-linux-amd64 -L -s -o manifest-tool
	@chmod +x manifest-tool
	@./manifest-tool push from-args \
		--platforms linux/arm,linux/arm64,linux/amd64 \
		--template "$(DOCKER_IMAGE):$(VERSION)-ARCH" \
		--target "$(DOCKER_IMAGE):$(VERSION)"

.PHONY: ci
ci: build test

.PHONY: test
test:
	@$(GO) test ./...

.PHONY: clean
clean:
	@rm -rf dist/
