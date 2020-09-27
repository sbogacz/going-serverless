SHELL=/bin/bash -o pipefail

GO_VERSION ?= 1.15
NAMESPACE=sbogacz
APPNAME=going-serverless
TALK_PHASE=02_separate_binaries
BINARY ?= handler

all: build

.PHONY: build
build: build-lambda

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

.PHONY: build-lambda
build-lambda:
	GOOS=linux go build -o $(BINARY) \
		github.com/$(NAMESPACE)/$(APPNAME)/$(TALK_PHASE)/cmd/lambda
	zip aws-$(BINARY).zip $(BINARY)

.PHONY: clean
clean: 
	go clean
	rm -rf $(BINARY) *.zip 
