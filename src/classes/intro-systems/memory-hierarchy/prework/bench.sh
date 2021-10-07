#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

./build.sh

hyperfine ./loop-order-one.out ./loop-order-two.out

valgrind --tool=cachegrind --cachegrind-out-file=cachegrind.out.loop-one ./loop-order-one.out
valgrind --tool=cachegrind --cachegrind-out-file=cachegrind.out.loop-two ./loop-order-two.out
