GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=k8s-yaml-splitter
TAG=latest

all: build

tidy:
	@go mod tidy

vendor:
	@go mod vendor

build: vendor
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOARCH=amd64 GOOS=linux   go build -o bin/$(GONAME)-linux $(GOFILES)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOARCH=amd64 GOOS=freebsd go build -o bin/$(GONAME)-freebsd $(GOFILES)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOARCH=amd64 GOOS=darwin  go build -o bin/$(GONAME)-darwin $(GOFILES)

container:
	@echo "Building container image"
	docker build -t ${GONAME}:${TAG} .

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	rm -rf ./bin
	rm -rf ./vendor

.PHONY: build vendor clean container
