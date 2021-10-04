#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

# function my_gcc() {
#   if [[ -f "/usr/local/bin/gcc-11" ]]; then
#     "/usr/local/bin/gcc-11" "$@"
#     return
#   fi

#   gcc "$@"
# }

# my_gcc -Og -S mstore.c
# my_gcc -Og -c mstore.c
# my_gcc -Og -o prog main.c mstore.c

./build.sh

hyperfine -i --warmup=500 --min-runs=1000 'pagecount_orig.o' 'pagecount_ctz.o' 'pagecount_encoding.o'
