#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"/..
set -eu pipefail

BUILD_FILES=()
mapfile -t BUILD_FILES < <(scripts/ls-build-files.sh)

./scripts/parallel_run.sh bash {} ::: "${BUILD_FILES[@]}"
