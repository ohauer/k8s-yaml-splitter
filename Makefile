GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=k8s-yaml-splitter
TAG=latest

help:
	@echo "Available targets:"
	@echo "  build      - Build binaries for all platforms"
	@echo "  test       - Run all tests with the built binary"
	@echo "  fmt        - Format Go source code"
	@echo "  container  - Build container image and extract binaries"
	@echo "  tidy       - Run go mod tidy"
	@echo "  vendor     - Download dependencies to vendor/"
	@echo "  clean      - Clean build artifacts"
	@echo "  help       - Show this help message"

fmt:
	@echo "Formatting Go code..."
	@go fmt ./...

tidy:
	@go mod tidy

vendor:
	@go mod vendor

build: vendor
	@echo "Building $(GOFILES) to ./bin"
	@mkdir -p bin
	@GOARCH=amd64 GOOS=linux   go build -o bin/$(GONAME)-linux-amd64 .
	@GOARCH=amd64 GOOS=freebsd go build -o bin/$(GONAME)-freebsd-amd64 .
	@GOARCH=amd64 GOOS=darwin  go build -o bin/$(GONAME)-darwin-amd64 .

test: build
	@echo "Running tests..."
	@./test.sh

container:
	@echo "Building container image"
	docker build -t ${GONAME}:${TAG} .

clean:
	@echo "Cleaning"
	@go clean
	rm -rf ./bin
	rm -rf ./vendor
	rm -rf ./test-output-*
	rm -f checksums.txt

.PHONY: build vendor clean container tidy fmt help test
.DEFAULT_GOAL := help
