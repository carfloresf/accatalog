#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

TRAVISBUILD=${TRAVIS:-}

if [ ! -z "${TRAVISBUILD}" ]; then
  echo "Updating golang dependencies"
  go get -u github.com/axw/gocov/gocov
  go get -u github.com/AlekSi/gocov-xml
  go get -u github.com/jstemmer/go-junit-report
  curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.15.0
fi

export GOROOT=$(go env GOROOT)
FILE="${FILE:-0}"

echo
echo "Running golangci"
	golangci-lint run --enable-all --disable=lll --disable=gocyclo --disable=gochecknoglobals --disable=gosimple --deadline=300s --tests=false ./...


echo
export CGO_ENABLED=0
echo "Running tests:"
  go test -v ./...


echo
echo "Testing coverage"

go test -covermode=count -coverprofile=profile.cov ./...
go tool cover -func profile.cov
  rm -f profile.cov
