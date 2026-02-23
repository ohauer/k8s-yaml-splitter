GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=k8s-yaml-splitter
TAG=latest
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -extldflags '-static' -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

help:
	@echo "Available targets:"
	@echo "  build        - Build binaries for all platforms"
	@echo "  test         - Run integration tests with the built binary"
	@echo "  test-unit    - Run unit tests with race detector"
	@echo "  test-all     - Run all tests (unit + integration)"
	@echo "  coverage     - Generate test coverage report"
	@echo "  fuzz         - Run fuzz tests (30s each)"
	@echo "  fmt          - Format Go source code"
	@echo "  lint         - Run linter (requires golangci-lint)"
	@echo "  container    - Build container image and extract binaries"
	@echo "  tidy         - Run go mod tidy"
	@echo "  vendor       - Download dependencies to vendor/"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"

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
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux   go build -a -trimpath $(LDFLAGS) -o bin/$(GONAME)-linux-amd64 .
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=freebsd go build -a -trimpath $(LDFLAGS) -o bin/$(GONAME)-freebsd-amd64 .
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin  go build -a -trimpath $(LDFLAGS) -o bin/$(GONAME)-darwin-amd64 .

test: build
	@echo "Running integration tests..."
	@./test.sh

test-unit:
	@echo "Running unit tests with race detector..."
	@go test -v -race ./...

test-all: test-unit test
	@echo "All tests completed successfully"

coverage:
	@echo "Generating test coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep total
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

fuzz:
	@echo "Running fuzz tests (30s each)..."
	@go test -fuzz=FuzzSanitizeFilename -fuzztime=30s
	@go test -fuzz=FuzzProcessObject -fuzztime=30s
	@go test -fuzz=FuzzParseFilterList -fuzztime=30s
	@echo "Fuzzing complete"

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run ./...

container:
	@echo "Building container and extracting binaries"
	@mkdir -p bin
	docker build --target builder -t ${GONAME}-builder:${TAG} .
	docker create --name temp-${GONAME}-builder ${GONAME}-builder:${TAG}
	docker cp temp-${GONAME}-builder:/go/src/github.com/ohauer/k8s-yaml-splitter/bin/. ./bin/
	docker rm temp-${GONAME}-builder
	docker build -t ${GONAME}:${TAG} .

clean:
	@echo "Cleaning"
	@go clean
	rm -rf ./bin
	rm -rf ./vendor
	rm -rf ./test-output-*
	rm -f checksums.txt
	rm -f coverage.out coverage.html

.PHONY: build vendor clean container tidy fmt help test test-unit test-all coverage fuzz lint
.DEFAULT_GOAL := help
