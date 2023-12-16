#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

prfx="[docs:lint]"

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)markdownlint"
NODE_NO_WARNINGS=1 markdownlint ./**/*.md

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)vale"
vale ./**/*.md
