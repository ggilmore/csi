#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

make clean
make

objdump -x86-asm-syntax=intel -S dotproduct.o >./dotproduct.s
