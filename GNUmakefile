default: build

install: vendor
	go install

build: vendor
	go build

vendor:
	dep ensure
	go generate

.PHONY: build vendor install
