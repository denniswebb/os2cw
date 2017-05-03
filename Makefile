OUTPUT_FILTER := build/{{.OS}}_{{.Arch}}/{{.Dir}}
CURRENT_DIR := $(shell pwd)
PROJECT := $(notdir $(CURRENT_DIR))
USER := $(notdir $(shell dirname $(CURRENT_DIR)))
CONTAINER_DIR := /go/src/github.com/$(USER)/$(PROJECT)
CONTAINER_DIR_CIRCLE := /go/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}
CIRCLECI := ${CIRCLECI}
GLIDE := $(shell glide -v dot 2> /dev/null)
GOX := $(shell gox -verbose dot 2> /dev/null)

.DEFAULT_GOAL := all

ifndef VERSION
  VERSION := git-$(shell git rev-parse --short HEAD)
endif

vend:
ifndef GLIDE
	$(shell curl https://glide.sh/get | sh)
endif
	glide install

fmt:
	goimports -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

all: vend
	gox  -ldflags "-X main.BuildVersion=${VERSION}" \
		-osarch darwin/amd64 -osarch linux/amd64 -osarch windows/amd64 \
		-output="$(OUTPUT_FILTER)"

linux: vend
	gox -ldflags "-X main.BuildVersion=${VERSION}" -osarch linux/amd64 -output="$(OUTPUT_FILTER)"

mac: vend
	gox -ldflags "-X main.BuildVersion=${VERSION}" -osarch darwin/amd64 -output="$(OUTPUT_FILTER)"

windows: vend
	gox -ldflags "-X main.BuildVersion=${VERSION}" -osarch windows/amd64 -output="$(OUTPUT_FILTER)"

clean:
	rm -rf build/
	go clean

on-docker:
ifeq ($(CIRCLECI), true)
	docker run -ti -v $(CURRENT_DIR):$(CONTAINER_DIR_CIRCLE) golang:1.7 /bin/bash -c "cd $(CONTAINER_DIR_CIRCLE) && ./build.sh"
else
	docker run -ti -v $(CURRENT_DIR):$(CONTAINER_DIR) golang:1.7 /bin/bash -c "cd $(CONTAINER_DIR) && ./build.sh"
endif

image: artifact
ifeq ($(CIRCLECI), true)
	docker build --rm=false -t ${CIRCLE_PROJECT_REPONAME}:$(shell head -1 VERSION).${CIRCLE_BUILD_NUM} .
	docker tag -f ${CIRCLE_PROJECT_REPONAME}:$(shell head -1 VERSION).${CIRCLE_BUILD_NUM} ${CIRCLE_PROJECT_REPONAME}:latest
else
	docker build -t $(PROJECT):latest .
endif
