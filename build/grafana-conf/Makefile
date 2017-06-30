IMAGE_NAME = martinbaillie/grafana-conf
IMAGE_VERSION = latest

.PHONY: build push

release:: build push

push::
	@docker push $(IMAGE_NAME):$(IMAGE_VERSION)

build::
	@docker build --pull \
		-t $(IMAGE_NAME):$(IMAGE_VERSION) $(shell pwd)
