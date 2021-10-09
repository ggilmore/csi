#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

./build.sh

hyperfine ./loop-order-one.out ./loop-order-two.out

hyperfine --parameter-list size 2,4,8,16,32,64,128,256,512,1024,2048 "./matrix-multiply.out {size}"

valgrind --tool=cachegrind --cachegrind-out-file=cachegrind.out.loop-one ./loop-order-one.out
valgrind --tool=cachegrind --cachegrind-out-file=cachegrind.out.loop-two ./loop-order-two.out
valgrind --tool=cachegrind --cachegrind-out-file=cachegrind.out.matrix-multiply ./matrix-multiply.out 256
