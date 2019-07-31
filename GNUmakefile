default: build

install:
	GOFLAGS=-mod=vendor GO111MODULE=on go install

build:
	GOFLAGS=-mod=vendor GO111MODULE=on go build

tools:
	GO111MODULE=off  go get -u github.com/alvaroloes/enumer

vendor:
	go mod vendor
	GOFLAGS=-mod=vendor GO111MODULE=on go generate

.PHONY: build vendor install tools
