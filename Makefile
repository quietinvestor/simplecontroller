DOCKER_TAG ?= simplecontroller:0.1.0

.PHONY: build fmt test vet docker-build kind-test kind-delete

build: fmt vet | bin
	go build -o bin/simplecontroller ./cmd

fmt:
	go fmt ./...

test: fmt vet
	go test ./...

vet:
	go vet ./...

bin:
	mkdir -p bin

docker-build:
	docker build -t $(DOCKER_TAG) .

kind-test: docker-build
	kind create cluster
	kind load docker-image $(DOCKER_TAG)
	kubectl apply -f config/default

kind-delete:
	kind delete cluster
