#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

gobin="${GOBIN:-go}"
prfx="[ci:lint]"

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)actionlint"
"$gobin" run github.com/rhysd/actionlint/cmd/actionlint@latest
