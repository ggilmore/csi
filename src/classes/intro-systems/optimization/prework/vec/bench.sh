#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

./build.sh

hyperfine -i --warmup=200 './tests'
