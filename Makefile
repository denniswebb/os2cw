OUTPUT_FILTER := build/{{.OS}}_{{.Arch}}/{{.Dir}}
CURRENT_DIR := $(shell pwd)
PROJECT := $(notdir $(CURRENT_DIR))
USER := $(notdir $(shell dirname $(CURRENT_DIR)))
CONTAINER_DIR := /go/src/github.com/$(USER)/$(PROJECT)
CONTAINER_DIR_CIRCLE := /go/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}
CIRCLECI := ${CIRCLECI}
GLIDE := $(shell glide -v dot 2> /dev/null)

all:
	gox -osarch darwin/amd64 -osarch linux/amd64 -output="$(OUTPUT_FILTER)"

linux:
	gox -osarch linux/amd64 -output="$(OUTPUT_FILTER)"

mac:
	gox -osarch darwin/amd64 -output="$(OUTPUT_FILTER)"

vend:
ifndef GLIDE
	$(shell curl https://glide.sh/get | sh)
endif
	glide install

clean:
	rm -rf build/
	go clean

artifact: clean
ifeq ($(CIRCLECI), true)
	docker run -ti -v $(CURRENT_DIR):$(CONTAINER_DIR_CIRCLE) golang:1.7 /bin/bash -c "cd $(CONTAINER_DIR_CIRCLE) && ./build.sh"
else
	docker run -ti -v $(CURRENT_DIR):$(CONTAINER_DIR) golang:1.7 /bin/bash -c "cd $(CONTAINER_DIR) && ./build.sh"
endif

image: artifact
ifeq ($(CIRCLECI), true)
	docker build --rm=false -t upsidetravel-docker.jfrog.io/${CIRCLE_PROJECT_REPONAME}:$(shell head -1 VERSION).${CIRCLE_BUILD_NUM} .
	docker tag -f upsidetravel-docker.jfrog.io/${CIRCLE_PROJECT_REPONAME}:$(shell head -1 VERSION).${CIRCLE_BUILD_NUM} upsidetravel-docker.jfrog.io/${CIRCLE_PROJECT_REPONAME}:latest
else
	docker build -t upsidetravel-docker.jfrog.io/$(PROJECT):latest .
endif
