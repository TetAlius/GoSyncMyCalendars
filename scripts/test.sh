#!/usr/bin/env bash

# scripts/test: Run test suite for application. Optionally pass in a path to an
#              individual test file to run a single test.

set -e
cd "$(dirname "$0")/.."

go test -race -coverprofile=coverage.txt -covermode=atomic ./...
#echo "" > coverage.txt

#for d in $(go list ./... | grep -v vendor); do
#    go test -race -coverprofile=profile.out -covermode=atomic $d
#    if [ -f profile.out ]; then
#        cat profile.out >> coverage.txt
#        rm profile.out
#    fi
#done
