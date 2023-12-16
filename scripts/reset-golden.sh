#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

dirs="$(find . -type d -name 'golden' -path '*/testdata/*')"

if [[ -z "$dirs" ]]; then
    echo "no golden test directories found"
    exit 0
fi

printf "> remove %d golden test directories? [y/N] " "${#dirs[@]}"
read -r -s -n 1 answer

if [[ "$answer" != "y" ]]; then
    echo
    echo "aborting"
    exit 0
fi

printf "\n\n"

for dir in "${dirs[@]}"; do
    rm -rf "$dir"
    echo "removed $dir"
done

