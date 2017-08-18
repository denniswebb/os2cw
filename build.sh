#!/usr/bin/env bash

set -ex

curl https://glide.sh/get | sh

go get github.com/mitchellh/gox
go install github.com/mitchellh/gox

make clean
make all
