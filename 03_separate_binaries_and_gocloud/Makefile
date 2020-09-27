SHELL=/bin/bash -o pipefail

GO_VERSION ?= 1.15
NAMESPACE=sbogacz
APPNAME=going-serverless
TALK_PHASE=03_separate_binaries_and_gocloud
BINARY ?= handler

all: build

.PHONY: build
build: build-aws build-gcp

.PHONY: build-local
build-local:
	go build -o $(BINARY) \
		github.com/$(NAMESPACE)/$(APPNAME)/$(TALK_PHASE)/cmd/http

.PHONY: build-gcp
build-gcp:
	mkdir tmp
	cp cmd/gcp/main.go tmp/function.go
	cp -R internal tmp/
	cp go.mod tmp/
	cp go.sum tmp/
	cd tmp && zip -r gcp-$(BINARY).zip . && cp *.zip ../ && cd -
	rm -rf tmp

.PHONY: build-aws
build-aws:
	GOOS=linux go build -o $(BINARY) \
		github.com/$(NAMESPACE)/$(APPNAME)/$(TALK_PHASE)/cmd/lambda
	zip aws-$(BINARY).zip $(BINARY)

.PHONY: clean
clean: 
	go clean
	rm -rf $(BINARY) *.zip 
