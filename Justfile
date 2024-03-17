set quiet := true
set dotenv-load := true

BIN_PATH := 'bin/lighthouse'
IMAGE_NAME := 'lighthouse_app'
REMOTE_DIR := 'lighthouse'
DB_NAME := 'lighthouse'
AUTOCAP_BIN_PATH := 'bin/autocap-$(hostname -s)'

mod make

[private]
help:
    just --list --unsorted

# debug disabled, services enabled
dev: build-autocap
    hivemind

# debug enabled, services disabled
dev2: build-autocap
    DISABLE_SERVICES=1 LIGHTHOUSE_DEBUG=1 hivemind

# debug enabled, services enabled
dev3: build-autocap
    LIGHTHOUSE_DEBUG=1 hivemind

db:
    mariadb {{ DB_NAME }}

drop:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

pull:
    importdb -f --host s --local {{ DB_NAME }} --remote {{ DB_NAME }}

fmt:
    go fmt ./...
    just --fmt

test: generate
    #!/bin/bash
    set -eo pipefail
    export VIEWS_DIR=$(realpath cmd/*/views)
    unbuffer go test -cover ./... | gostack --test

cov: generate
    #!/bin/bash
    set -eo pipefail
    FILE=$(mktemp)
    export VIEWS_DIR=$(realpath cmd/*/views)
    go test ./... -cover -coverprofile="$FILE"
    go tool cover -html="$FILE"
    rm "$FILE"

lint:
    #!/bin/bash
    set -eo pipefail
    unbuffer go vet ./... | gostack
    unbuffer golangci-lint --color never run | gostack

fix:
    #!/bin/bash
    set -eo pipefail
    unbuffer golangci-lint --color never run --fix | gostack

sloc:
    tokei -tGo,HTML,CSS,JavaScript \
        -e **/alpine.min.js \
        -e **/bulma.min.css

clean:
    rm -vf {{ AUTOCAP_BIN_PATH }}
    rm -vf {{ BIN_PATH }}

deploy: buildc
    deploy-container \
        --host s \
        --image {{ IMAGE_NAME }} \
        --dir {{ REMOTE_DIR }} \
        --compose docker-compose.yml \
        --add .env.prod \
        --init deploy/init

buildc:
    docker-compose build

# ——————————————————————————————————————————————————————————————————————————————————————————————————

[private]
up: buildc
    docker-compose up

[private]
down:
    docker-compose down

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
    #!/bin/bash
    set -eo pipefail
    {{ AUTOCAP_BIN_PATH }} {{ BIN_PATH }}
    unbuffer {{ BIN_PATH }} | gostack

[private]
generate:
    go generate ./...

[private]
build: generate
    #!/bin/bash
    set -eo pipefail
    unbuffer go build -o {{ BIN_PATH }} ./cmd/lighthouse | gostack
