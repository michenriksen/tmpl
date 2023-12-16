#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

gobin="${GOBIN:-go}"
goversion="$("$gobin" env GOVERSION)"

export GO_VERSION="$goversion"
export GITLAB_TOKEN=""

goreleaser build --clean --snapshot --single-target
