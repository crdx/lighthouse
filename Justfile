set quiet

BIN_PATH := 'dist/lighthouse'
AUTOCAP_BIN_PATH := 'bin/autocap-$(hostname -s)'

mod make
mod container
mod db

set dotenv-load := true

[private]
help:
    just --list --unsorted

# start development
dev: build-autocap
    watchexec -r 'just redev'

# start development with debug logging and no services
devd:
    DISABLE_SERVICES=1 LIGHTHOUSE_DEBUG=1 just dev

# start development with debug logging
devdd:
    LIGHTHOUSE_DEBUG=1 just dev

# build binary
build:
    mkif {{ BIN_PATH }} $(find src -type f) -x 'just rebuild'

# run tests
test:
    #!/bin/bash
    set -eo pipefail
    cd src
    export VIEWS_DIR=$(realpath views)
    go test -cover ./... 2>&1 | gostack --mode gotest --mod crdx.org/lighthouse
    exit ${PIPESTATUS[0]}

# run tests and show code coverage
cov:
    #!/bin/bash
    set -e
    cd src
    FILE=$(mktemp)
    export VIEWS_DIR=$(realpath views)
    go test -cover -coverprofile="$FILE" ./...
    go tool cover -html="$FILE"
    rm "$FILE"

# run a specific test
test-file path:
    #!/bin/bash
    set -e
    cd src
    FILE="{{ path }}"
    export VIEWS_DIR=$(realpath views)
    go test -cover "${FILE#src/}"
    # go test -cover -cpuprofile ../cpu.prof "${FILE#src/}"
    # go tool pprof -png -output ../cpu.png ../cpu.prof

# check everything
check: lint test

# run linter
lint:
    # We don't need the loopclosure check because of GOEXPERIMENT=loopvar.
    cd src && go vet -loopclosure=false ./... && golangci-lint run

# run formatter
fmt:
    cd src && go fmt ./...

# show code stats
sloc:
    tokei -tGo,HTML,CSS,JavaScript \
        -e src/assets/alpine.min.js \
        -e src/assets/bulma.min.css

# deploy the container
deploy:
    just container::deploy

# ——————————————————————————————————————————————————————————————————————————————————————————————————

[private]
build-autocap:
    #!/bin/bash
    mkif {{ AUTOCAP_BIN_PATH }} src/cmd/autocap/main.go -x 'just rebuild-autocap'

[private]
rebuild-autocap:
    #!/bin/bash
    set -e
    sudo rm -f {{ AUTOCAP_BIN_PATH }}
    go build -o {{ AUTOCAP_BIN_PATH }} src/cmd/autocap/main.go
    sudo chown root {{ AUTOCAP_BIN_PATH }}
    sudo chmod u+s {{ AUTOCAP_BIN_PATH }}

[private]
redev: build
    #!/bin/bash
    set -eo pipefail
    {{ AUTOCAP_BIN_PATH }} {{ BIN_PATH }}
    unbuffer {{ BIN_PATH }} --env .env 2>&1 | gostack --mod crdx.org/lighthouse

[private]
generate:
    cd src && go generate ./...

[private]
rebuild: generate
    #!/bin/bash
    set -eo pipefail
    cd src
    unbuffer go build -o ../{{ BIN_PATH }} -trimpath -ldflags '-s -w' 2>&1 | gostack --mod crdx.org/lighthouse

[private]
clean:
    rm -fv {{ BIN_PATH }}
