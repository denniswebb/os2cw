GO_VERSION := latest
circleci := ${CIRCLECI}
output_filter := build/{{.OS}}_{{.Arch}}/{{.Dir}}
current_dir := $(shell pwd)
user := $(notdir $(shell dirname $(current_dir)))
project := $(notdir $(current_dir))
gitsha := $(shell git rev-parse HEAD)
build_date := $(shell date -Iseconds)
glide := $(shell glide -v dot 2> /dev/null)
container_dir := /go/src/github.com/$(user)/$(project)
container_dir_circle := /go/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}

.DEFAULT_GOAL := all

ifndef VERSION
  VERSION := git-$(shell git rev-parse --short HEAD)
endif

vend:
ifndef glide
	$(shell curl https://glide.sh/get | sh)
endif
	glide install

fmt:
	goimports -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

all: vend
	gox  -ldflags "-X main.BuildVersion=${VERSION}" \
		-osarch darwin/amd64 -osarch linux/amd64 -osarch windows/amd64 \
		-output="$(output_filter)"

linux: vend
	gox -ldflags "-X main.BuildVersion=${VERSION}" -osarch linux/amd64 -output="$(output_filter)"

mac: vend
	gox -ldflags "-X main.BuildVersion=${VERSION}" -osarch darwin/amd64 -output="$(output_filter)"

windows: vend
	gox -ldflags "-X main.BuildVersion=${VERSION}" -osarch windows/amd64 -output="$(output_filter)"

clean:
	rm -rf build/
	go clean

on-docker:
	@echo $(container_dir)
ifeq ($(CIRCLECI), true)
	docker create -v $(container_dir_circle) --name src alpine:3.4 /bin/true
	docker cp $(current_dir)/. src:$(container_dir_circle)
	docker run -ti --volumes-from src golang:$(GO_VERSION) /bin/bash -c "cd $(container_dir_circle) && ./build.sh"
	docker cp src:$(container_dir_circle)/build/. $(current_dir)/build
else
	docker run -ti -v $(current_dir):$(container_dir) golang:$(GO_VERSION) /bin/bash -c "cd $(container_dir) && ./build.sh"
endif
