SHELL=/bin/bash -o pipefail

GO_VERSION ?= 1.15
NAMESPACE=sbogacz
APPNAME=going-serverless
TALK_PHASE=04_build_tags
BINARY ?= handler

all: build

.PHONY: build
build: build-local build-aws

.PHONY: build-local
build-local:
	go build -o $(BINARY) \
		github.com/$(NAMESPACE)/$(APPNAME)/$(TALK_PHASE)/cmd/toy


.PHONY: build-aws
build-aws:
	GOOS=linux go build -tags aws -o $(BINARY) \
		github.com/$(NAMESPACE)/$(APPNAME)/$(TALK_PHASE)/cmd/toy
	zip aws-$(BINARY).zip $(BINARY)

.PHONY: clean
clean: 
	go clean
	rm -rf $(BINARY) *.zip 
