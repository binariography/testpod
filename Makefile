
TAG?=latest
NAME:=testpod
DOCKER_REPOSITORY:=ciavash
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
VERSION=$(shell grep 'VERSION' pkg/version/version.go | awk '{print $$4}' | tr -d '"')
GO:=go1.25.4

test: tidy fmt vet
	$(GO) test ./...

build:
	$(GO) build -ldflags "-s -w" -o ./bin/$(NAME)

tidy:
	rm -f go.sum; $(GO) mod tidy -compat=1.25

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...


docker_build:
	docker build --tag $(DOCKER_IMAGE_NAME):$(VERSION) .

image_push:
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):$(TAG)
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest

