set quiet := true
set shell := ["bash", "-cu", "-o", "pipefail"]
set dotenv-load := true

NAME := 'lighthouse'
IMAGE_NAME := 'lighthouse_app'
REMOTE_DIR := 'lighthouse'
DB_NAME := 'lighthouse'
HOST := 'm'
BIN_PATH := 'bin/' + NAME
AUTOCAP_BIN_PATH := 'bin/autocap-$(hostname -s)'

import? 'local.just'

mod make

[private]
help:
    just --list --unsorted --list-submodules

# debug enabled, services disabled
dev: build-autocap
    LIGHTHOUSE_DEBUG=1 DISABLE_SERVICES=1 hivemind

# debug disabled, services enabled
dev2: build-autocap
    hivemind

# debug enabled, services enabled
dev3: build-autocap
    LIGHTHOUSE_DEBUG=1 hivemind

db:
    mariadb {{ DB_NAME }}

drop:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

fetchdb:
    importdb -f --host {{ HOST }} --local {{ DB_NAME }} --remote {{ DB_NAME }}

fmt:
    go fmt ./...

test package='./...': generate
    LOG_TYPE=none unbuffer go test -cover {{ package }} | gostack --test

cov package='./...': (gencov 'func' package)
cov-html package='./...': (gencov 'html' package)

[private]
gencov flag package: generate
    #!/bin/bash
    set -euo pipefail
    FILE=$(mktemp)
    export LOG_TYPE=none
    unbuffer go test {{ package }} -cover -coverprofile="$FILE" | gostack --test
    sed -i '/\.gen\.go/d' "$FILE"
    if [ "{{ flag }}" = func ]; then
        unbuffer go tool cover -{{ flag }}="$FILE" | gostack --test
    else
        go tool cover -{{ flag }}="$FILE"
    fi
    rm "$FILE"

lint:
    unbuffer go vet ./... | gostack
    unbuffer golangci-lint --color never run | gostack

fix:
    unbuffer golangci-lint --color never run --fix | gostack

deploy: buildc
    deploy-container \
        --host {{ HOST }} \
        --image {{ IMAGE_NAME }} \
        --dir {{ REMOTE_DIR }} \
        --compose compose.yml \
        --add .env.prod \
        --init deploy/init

check: test fmt lint

build: generate
    unbuffer go build -o {{ BIN_PATH }} ./cmd/{{ NAME }} | gostack

buildc: generate
    docker-compose build

clean:
    rm -vf {{ AUTOCAP_BIN_PATH }}
    rm -vf {{ BIN_PATH }}

[private]
build-autocap:
    mkif {{ AUTOCAP_BIN_PATH }} ./tools/autocap/* -x 'just rebuild-autocap'

[private]
rebuild-autocap:
    sudo rm -f {{ AUTOCAP_BIN_PATH }}
    go build -o {{ AUTOCAP_BIN_PATH }} ./tools/autocap
    sudo chown root {{ AUTOCAP_BIN_PATH }}
    sudo chmod u+s {{ AUTOCAP_BIN_PATH }}

[private]
serve: build
    {{ AUTOCAP_BIN_PATH }} {{ BIN_PATH }}
    unbuffer {{ BIN_PATH }} | gostack

[private]
generate:
    go generate ./...
