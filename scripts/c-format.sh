#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"/..
set -eu pipefail

C_FILES=()
mapfile -t C_FILES < <(scripts/ls-c-files.sh)

function format() {
  local file="$1"

  local FORMAT_ARGS=(
    "-i"
    "${file}"
  )

  result=$(clang-format "${FORMAT_ARGS[@]}" 2>&1)
  rc=$?

  if [ -n "$result" ]; then
    echo "${file}:"
    echo "$result"
    echo
  fi

  exit "$rc"
}
export -f format

./scripts/parallel_run.sh format {} ::: "${C_FILES[@]}"
