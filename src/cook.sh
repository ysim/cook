#!/usr/bin/env bash

set -o pipefail

RECIPES_DIR="${COOK_RECIPES_DIR:-$HOME/.cook-recipes}"
RECIPES_EXTENSION=".md"

die() {
    echo "$@" >&2
    exit 1
}

cmd_show() {
    local name="$1"
    local recipe_path="${RECIPES_DIR}/${name}${RECIPES_EXTENSION}"

    if [[ -f $recipe_path ]]; then
        cat "$recipe_path"
    else
        echo "No recipe exists at the path: ${recipe_path}"
    fi
}

cmd_other() {
    die "Not a valid subcommand."
}

COMMAND="$1"

case "${COMMAND}" in
    show) shift; cmd_show "$@" ;;
    *) cmd_other "$@" ;;
esac
exit 0
