all: check format lint build

build: render-ci-pipeline build-all

render-ci-pipeline:
    ./scripts/render-ci-pipeline.sh

fmt: format
format: format-dhall prettier format-shfmt gofmt format-c

lint: lint-dhall shellcheck

check: check-dhall

gen: generate
generate:
    ./ast_generator/generate.sh

build-all:
    ./scripts/build-all.sh

gofmt: format-golang
format-golang:
    ./scripts/go-fmt.sh

format-c:
    ./scripts/c-format.sh

test: test-golang

test-golang:
    ./scripts/go-test.sh

prettier:
    yarn run prettier

format-dhall:
    ./scripts/dhall-format.sh

check-dhall:
    ./scripts/dhall-check.sh

lint-dhall:
    ./scripts/dhall-lint.sh

shellcheck:
    ./scripts/shellcheck.sh

format-shfmt:
    shfmt -w .

install:
    just install-asdf
    just install-yarn

install-yarn:
    yarn

install-asdf:
    asdf install
