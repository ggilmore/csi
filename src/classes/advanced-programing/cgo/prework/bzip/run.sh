#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

go build -o bzipper ./...

wc -c </usr/share/dict/words
sha256sum </usr/share/dict/words
./bzipper </usr/share/dict/words | wc -c
./bzipper </usr/share/dict/words | bunzip2 | sha256sum
