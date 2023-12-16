#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

gobin="${GOBIN:-go}"
testflags=("-coverpkg=./..." "-covermode=atomic" "-coverprofile=coverage.out")

if [[ "${DBG:-}" == 1 ]]; then
    testflags+=("-v")
fi

printf "$(tput bold)[go:test] running %s...$(tput sgr0)\n" "$(tput setaf 4)go test"
"$gobin" test "${testflags[@]}" ./...

"$gobin" tool cover -html=coverage.out
