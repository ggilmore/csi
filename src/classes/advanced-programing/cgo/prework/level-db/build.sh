#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

# main() {
#   gcc "-I$PWD/leveldb/include/" -I/usr/include "-L$PWD/leveldb/build" \
# 	-lstdc++ -l leveldb -lleveldb-dev \
#     main.c
# }

# main "$@"

rm -rf "abc" || true 
go run . 