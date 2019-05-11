export CGO_ENABLED = 0

GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
ARCH := $(GOARCH)
ifeq ($(ARCH), arm)
	ARCH = armhf
endif

PROTO_TOOL = docker run --rm -v `pwd`:/work uber/prototool:1.7.0 prototool
PROTO_PACKAGE := api
PROTO_FILES := $(wildcard $(PROTO_PACKAGE)/*.proto)
GENERATED_CODE := $(patsubst %.proto, %.pb.go, $(PROTO_FILES))

DOCKER_IMAGE := jjbubudi/tides
TAG ?= $(shell git describe --tags 2>/dev/null)
COMMIT ?= $(shell git rev-parse --short HEAD)
VERSION := $(COMMIT)
ifneq ($(TAG),)
	VERSION := $(TAG)
endif

.PHONY: build
build: generate
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
		--target "$(DOCKER_IMAGE):latest"
	@./manifest-tool push from-args \
		--platforms linux/arm,linux/arm64,linux/amd64 \
		--template "$(DOCKER_IMAGE):$(VERSION)-ARCH" \
		--target "$(DOCKER_IMAGE):$(VERSION)"

.PHONY: run
run: generate
	@$(GO) run cmd/main.go

.PHONY: ci
ci: build test

.PHONY: test
test: generate
	@$(GO) test ./...

.PHONY: generate
generate: $(GENERATED_CODE)

%.pb.go: %.proto
	@$(PROTO_TOOL) generate

.PHONY: clean
clean:
	@rm -rf dist/
	@rm -rf api/*.pb.go