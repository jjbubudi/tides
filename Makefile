GO = go
PROTO_TOOL = docker run --rm -v `pwd`:/work uber/prototool:1.7.0 prototool
PROTO_PACKAGE = api
PROTO_FILES = $(wildcard $(PROTO_PACKAGE)/*.proto)
GENERATED_CODE = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))

.PHONY: build
build: generate
	@CGOENABLED=0 $(GO) build -o dist/tides cmd/main.go

.PHONY: run
run: generate
	@$(GO) run cmd/main.go

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