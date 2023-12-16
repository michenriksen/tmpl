#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

printf "$(tput bold)[yaml:lint] running %s...$(tput sgr0)\n" "$(tput setaf 4)yamllint"
yamllint -s -f colored .
