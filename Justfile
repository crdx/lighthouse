BIN_PATH := 'dist/lighthouse'
IMAGE_NAME := 'lighthouse_app'
REMOTE_DIR := 'lighthouse'
DB_NAME := 'lighthouse'
AUTOCAP_BIN_PATH := 'bin/autocap-$(hostname -s)'

set dotenv-load := true

[private]
@help:
    just --list --unsorted

# start development
@dev: make-autocap
    watchexec \
        --no-vcs-ignore \
        -w .env \
        -w src \
        -i src/conf/models.go \
        -i src/conf/migrations.go \
        -r 'just redev'

# start development without live reload
@devn: make-autocap
    LIVE_RELOAD=0 just dev

# start development with debug logging and no services
@devd:
    DISABLE_SERVICES=1 LIGHTHOUSE_DEBUG=1 just dev

# start development with debug logging
@devdd:
    LIGHTHOUSE_DEBUG=1 just dev

# open the development site in the default browser
@open:
    xdg-open http://$HOST:$PORT

# build binary
@make:
    mkif {{ BIN_PATH }} $(find src -type f) -x 'just remake'

# build the container
@build:
    docker-compose build

# start the container
@up: build
    docker-compose up

# stop the container
@down:
    docker-compose down

# deploy the container
@deploy: build
    deploy-container \
        --host s \
        --image "{{ IMAGE_NAME }}" \
        --dir "{{ REMOTE_DIR }}" \
        --compose docker-compose.yml \
        --add .env.prod \
        --init deploy/init

# start a shell in the app container
@shell: build
    docker-compose run app bash || true
    docker-compose down

# run tests
test:
    #!/bin/bash
    set -e
    cd src
    export VIEWS_DIR=$(realpath views)
    go test -cover ./... | grep -vF '[no test files]'
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
@check: lint test

# run linter
@lint:
    # We don't need the loopclosure check because of GOEXPERIMENT=loopvar.
    cd src && go vet -loopclosure=false ./... && golangci-lint run

# run formatter
@fmt:
    cd src && go fmt ./...

# show code stats
@sloc:
    tokei -tGo,HTML,CSS,JavaScript \
        -e src/assets/alpine.min.js \
        -e src/assets/bulma.min.css

# connect to the dev db
@db:
    mariadb {{ DB_NAME }}

# drop the dev db
@drop-db:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

# connect to the live db
@live-db:
    ssh-to s -t mariadb {{ DB_NAME }}

# pull down the live db
@pull:
    importdb -f --host s --local {{ DB_NAME }} --remote {{ DB_NAME }}

# scaffold a new service (name is lowercase e.g. watcher)
@new-service name:
    mkdir -pv src/services/{{ name }}
    touch src/services/{{ name }}/{{ name }}.go
    echo 'package {{ name }}' > src/services/{{ name }}/{{ name }}.go

# scaffold a new model (name is lowercase e.g. device)
@new-model name:
    touch src/m/{{ name }}.go
    echo 'package m' > src/m/{{ name }}.go
    mkdir src/m/repo/{{ name }}R
    touch src/m/repo/{{ name }}R/{{ name }}.go
    echo 'package {{ name }}R' > src/m/repo/{{ name }}R/{{ name }}.go

# scaffold a new controller (name is lowercase e.g. device)
@new-controller name:
    mkdir -pv src/controllers/{{ name }}Controller
    touch src/controllers/{{ name }}Controller/routes.go
    touch src/controllers/{{ name }}Controller/routes_test.go
    echo 'package {{ name }}Controller' > src/controllers/{{ name }}Controller/routes.go
    echo 'package {{ name }}Controller_test' > src/controllers/{{ name }}Controller/routes_test.go

# scaffold a new migration (name is ProperCase with no punctuation)
[no-exit-message]
@new-migration name:
    go run src/cmd/mkmigration/main.go "{{ name }}"

[private]
make-autocap:
    #!/bin/bash
    mkif {{ AUTOCAP_BIN_PATH }} src/cmd/autocap/main.go -x 'just remake-autocap'

[private]
remake-autocap:
    #!/bin/bash
    set -e
    sudo rm -f {{ AUTOCAP_BIN_PATH }}
    go build -o {{ AUTOCAP_BIN_PATH }} src/cmd/autocap/main.go
    sudo chown root {{ AUTOCAP_BIN_PATH }}
    sudo chmod u+s {{ AUTOCAP_BIN_PATH }}

[private]
@redev: make
    {{ AUTOCAP_BIN_PATH }} {{ BIN_PATH }}
    {{ BIN_PATH }} --env .env

[private]
@generate:
    cd src && go generate ./...

[private]
@remake: generate
    cd src && go build -o ../{{ BIN_PATH }} -trimpath -ldflags '-s -w'

[private]
@clean:
    rm -fv {{ BIN_PATH }}
