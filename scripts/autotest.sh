#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

watchexec       \
  -c clear      \
  -o do-nothing \
  -d 100ms      \
  --exts go     \
  --shell=bash  \
  'pkg=".${WATCHEXEC_COMMON_PATH/$PWD/}/..."; echo "running tests for $pkg"; go test "$pkg"'
