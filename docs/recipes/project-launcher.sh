#!/usr/bin/env bash

set -euo pipefail

style_boldblue="\033[1;34m"
style_reset="\033[0m"

# if $NO_COLOR is set, set the style variables to empty strings.
if [ "${NO_COLOR:-}" != "" ]; then
  style_boldblue=""
  style_reset=""
fi

## Absolute path to directory where projects are stored.
#
# Change this to point to the directory where you store your projects.
projects_dir="$HOME/projects"

## Match pattern for identifying project root directories.
#
# This pattern finds directories that contain a .git directory. You can change
# this to match something else if you don't use git for your projects.
match_pattern='^\.git$'

## Path patterns to exclude from the search.
#
# Add any directories that you want to exclude from the search to this array.
exclude_patterns=(
  "node_modules"
  "vendor"
)

## Arguments to pass to fd.
#
# See `man fd` for more information.
fd_args=(
  "-H"        # Include hidden files and directories. This is needed to match .git
  "-t" "d"    # Only match directories. Change to "-t" "f" to match files instead.
  "--prune"   # Don't traverse matching directories.
)
for pattern in "${exclude_patterns[@]}"; do
  fd_args+=("-E" "$pattern")
done

## Use caching for project directories.
#
# Caches project directories to a file for faster startup time. Set to false if
# you don't want to use caching.
#
# You can clear the cache by setting $CLEAR_CACHE to any value when running
# the project launcher.
#
# You can also temporarily disable the cache by setting $NO_CACHE to any value
# when running the project launcher.
use_cache=true
if [ "${NO_CACHE:-}" != "" ]; then
  use_cache=false
fi

## Cache file location.
#
# Where the cache file is stored. Only used if $use_cache is true.
cache_file="${XDG_CACHE_HOME:-$HOME/Library/Caches}/tmpo/projects.cache"

## Get project data.
#
# Function returns project data to be used by fzf for the project selection
# menu.
#
# Caching of data is also managed by this function.
function get_project_data() {
  # Clear the cache if $CLEAR_CACHE is set.
  if [ "${CLEAR_CACHE:-}" != "" ] && [ -f "$cache_file" ]; then
    rm "$cache_file"
  fi

  # Return the cached data if caching is enabled and the cache file exists.
  if "$use_cache" && [[ -f "$cache_file" ]]; then
    cat "$cache_file"
    return
  fi

  project_data=""
  projects="$(fd "${fd_args[@]}" "$match_pattern" "$projects_dir" | xargs dirname)"

  while IFS= read -r path; do
    # Remove the repeditive projects directory prefix from the path to make it
    # more concise and readable.
    pretty_path="${path#"$projects_dir"/}"

    # To make the project selection more user-friendly, the project name is
    # extracted using the basename command and styled in bold and blue to make
    # it stand out.
    name="$style_boldblue$(basename "$pretty_path")$style_reset"
    pretty_path="$(dirname "$pretty_path")/$name"

    # The project path and its pretty path are appended to the data, separated
    # by a tab character ("\t"). This data will be used by fzf for the project
    # selection menu.
    project_data+="$path\t$pretty_path\n"
  done <<< "$projects"

  # Save the data to the cache file if caching is enabled.
  if "$use_cache"; then
    mkdir -p "$(dirname "$cache_file")"
    echo -e "$project_data" > "$cache_file"
  fi

  echo -e "$project_data"
}

# Present the selection menu using fzf and store the selected project path.
#
# Fzf is configured to split each line of project data by \t and use the first
# column for fuzzy matching and the second column for display.
selected_project="$(
  get_project_data |
    fzf \
      --delimiter="\t" \
      --nth=1 \
      --with-nth=2 \
      --scheme="path" \
      --no-info \
      --no-scrollbar \
      --ansi |
    cut -d $'\t' -f 1
)"

if [ "$selected_project" = "" ]; then
  # If no project was selected, exit the script.
  exit 0
fi

# Change directory to selected project and run tmpl.
(cd "$selected_project" || exit 1; tmpl)

