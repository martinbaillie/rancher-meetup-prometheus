SHELL 	:=/bin/bash
DIR 	:=$(shell pwd)
PROJECT	?=$(shell basename ${PWD})
VERSION ?=latest
TAG 	?=martinbaillie/$(PROJECT):$(VERSION)

release:: build push

push::
	@docker push $(TAG)

build::
	@docker build -t $(TAG) $(DIR)

.PHONY: build push release
