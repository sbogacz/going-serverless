SHELL=/bin/bash -o pipefail

GO_VERSION ?= 1.15
BINARY ?= handler
WD ?= $(shell pwd)
GO_FILES = $(shell find . -type f -name "*.go")

all: build

.PHONY: build
build: 
	GOOS=linux go build -o $(BINARY)
	zip $(BINARY).zip $(BINARY)

.PHONY: clean
clean: 
	go clean
	rm -rf $(BINARY) $(BINARY).zip
