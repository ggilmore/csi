#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

./build.sh

hyperfine -i --export-markdown=pagecount_timing.md './pagecount_orig.o' './pagecount_ctz.o' './pagecount_encoding.o'
