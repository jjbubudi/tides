export CGO_ENABLED = 0
export GOOS = $(shell go env GOOS)
export GOARCH = $(shell go env GOARCH)

.PHONY: build
build:
	@go build -o dist/tides cmd/main.go

.PHONY: ci
ci: build test

.PHONY: test
test:
	@go test ./...

.PHONY: clean
clean:
	@rm -rf dist/
