#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

OPTIMIZATION="-O0"

cc "${OPTIMIZATION}" -o loop-order-one.out loop-order-one.c loop-order.h
cc "${OPTIMIZATION}" -S loop-order-one.c loop-order.h

cc "${OPTIMIZATION}" -o loop-order-two.out loop-order-two.c loop-order.h
cc "${OPTIMIZATION}" -S loop-order-two.c loop-order.h
