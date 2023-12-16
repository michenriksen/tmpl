#!/usr/bin/env bash

set -e

if [[ "${DBG:-}" == 1 ]]; then
    set -x
fi

gobin="${GOBIN:-go}"

prfx="[docs:gen]"

schema_tmpdest="docs/jsonschema.md.tmp"
schema_dest="docs/jsonschema.md"

printf "$(tput bold)%s running %s...$(tput sgr0)\n" "$prfx" "$(tput setaf 4)go generate"
"$gobin" generate -tags=gen ./...

if [[ -f "$schema_dest" ]]; then
  printf "$(tput bold)$prfx removing %s...$(tput sgr0)\n" "$(tput setaf 4)$schema_dest"
  rm "$schema_dest"
fi

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)generate-schema-doc"
generate-schema-doc \
  --config template_name=md \
  --config link_to_reused_ref=false \
  --config show_toc=false \
  --config examples_as_yaml=true \
  --config footer_show_time=false \
  config.schema.json "$schema_tmpdest" > /dev/null

printf "$(tput bold)$prfx running %s...$(tput sgr0)\n" "$(tput setaf 4)mdformat"
mdformat "$schema_tmpdest"

printf "$(tput bold)%s removing HTML tags...$(tput sgr0)\n" "$prfx"
sed 's/<[^>]*>//g' "$schema_tmpdest" > "$schema_dest"
rm "$schema_tmpdest"
