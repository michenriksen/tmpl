#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

gobin="${GOBIN:-go}"
schema="config.schema.json"
config_files=("docs/.tmpl.reference.yaml")
prfx="[schema:lint]"

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)jv on $schema"
"$gobin" run github.com/santhosh-tekuri/jsonschema/cmd/jv@latest -assertcontent --assertformat "$schema"

for config_file in "${config_files[@]}"; do
    if [ ! -f "$config_file" ]; then
        printf "$(tput bold)$prfx %s not found$(tput sgr0)\n" "$config_file"
        exit 1
    fi

    printf "$(tput bold)$prfx running $(tput setaf 4)jv on %s...$(tput sgr0)\n" "$config_file"
    "$gobin" run github.com/santhosh-tekuri/jsonschema/cmd/jv@latest -assertcontent --assertformat config.schema.json "$config_file"
done
