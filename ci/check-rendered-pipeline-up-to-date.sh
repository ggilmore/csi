#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"/..
set -euxo pipefail

just render-ci-pipeline
just format

git diff --exit-code .github/workflows/ci.yaml
