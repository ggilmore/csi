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

cc -O1 -S -masm=intel pagecount_orig.c
cc -O1 -S -masm=intel pagecount_encoding.c
cc -O1 -S -masm=intel pagecount_ctz.c

cc -g3 -o pagecount_orig.o pagecount_orig.c
cc -g3 -o pagecount_encoding.o pagecount_encoding.c
cc -g3 -o pagecount_ctz.o pagecount_ctz.c
