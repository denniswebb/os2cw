#!/usr/bin/env bash

set -ex

go get -t -d -v $(go list ./... | grep -v /vendor/)
go test -v -race $(go list ./... | grep -v /vendor/)
mkdir -p build/
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o build/app
