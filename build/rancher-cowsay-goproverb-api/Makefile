DIR 	:=$(shell pwd)
TIME 	:=$(shell date '+%Y-%m-%dT%T%z')
HASH 	:=$(shell git rev-parse --verify --short HEAD || echo "")
USER 	:=$(shell echo $$(whoami))

PROJECT	?=$(shell basename ${PWD})
VERSION ?=latest
HOST 	?=github.com
TAG 	?=martinbaillie/$(PROJECT):$(VERSION)

LDFLAGS ?=-s -w -extld ld -extldflags -static \
		  -X 'main.buildTime=$(TIME)' \
		  -X 'main.buildHash=$(HASH)' \
		  -X 'main.buildUser=$(USER)' \
		  -X 'main.projectName=$(PROJECT)' \
		  -X 'main.projectVersion=$(VERSION)'

release:: build push

dep::
	@go get -u github.com/golang/dep/...

dep-update:: dep
	@dep ensure -update

push::
	@docker push $(TAG)

build::
	@docker build \
		--build-arg HOST='$(HOST)' \
		--build-arg PROJECT='$(PROJECT)' \
		--build-arg LDFLAGS='$(LDFLAGS)' \
		-t $(TAG) $(DIR)

.PHONY: build push dep-update dep release
