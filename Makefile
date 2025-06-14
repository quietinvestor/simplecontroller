DOCKER_TAG ?= simplecontroller:0.1.0

.PHONY: build test delete

build:
	docker build -t $(DOCKER_TAG) .

test: build
	kind create cluster
	kind load docker-image $(DOCKER_TAG)
	kubectl apply -f config/default

delete:
	kind delete cluster
