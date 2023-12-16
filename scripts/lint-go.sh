#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

gobin="${GOBIN:-go}"
prfx="[go:lint]"

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)go vet"
"$gobin" vet ./...

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)go mod verify"
"$gobin" mod verify > /dev/null

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)golangci-lint"
golangci-lint run
