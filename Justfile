BIN_PATH := 'dist/lighthouse'
IMAGE_NAME := 'lighthouse_app'
REMOTE_DIR := 'lighthouse'
DB_NAME := 'lighthouse'
AUTOCAP_BIN_PATH := 'bin/autocap-$(hostname -s)'

set dotenv-load := true

@_help:
    just --list --unsorted

# start development
@dev: make-autocap
    watchexec -w src -r 'just redev'

# start development with debug logging
@devd:
    LIGHTHOUSE_DEBUG=1 just dev

# start development with verbose logging
@devv:
    LIGHTHOUSE_VERBOSE=1 just dev

# build binary
@make:
    mkif {{ BIN_PATH }} $(find src -type f) -x 'just remake'

# remove build
@clean:
    rm -fv {{ BIN_PATH }}

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
    docker-compose run --service-ports app bash || true
    docker-compose down

# watch running container logs
@logs:
    docker-compose logs -f

# run tests
@test *args:
    cd src && go test -cover ./... {{ args }} | grep -vF '[no test files]' || true

# connect to the test database
@test-db:
    mariadb {{ DB_NAME }}_test

# run linter
@lint:
    cd src && golangci-lint run

# fix lint issues
@fix:
    cd src && golangci-lint run --fix

# drop the database
@drop-db:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

# connect to the dev database
@db:
    mariadb {{ DB_NAME }}

# pull down the live database
@pull:
    importdb -f --host s --local {{ DB_NAME }} --remote {{ DB_NAME }}

# scaffold a new model (name is lowercase e.g. thing)
@new-model name:
    mkdir -pv src/models/{{ name }}Model
    touch src/models/{{ name }}Model/{{ name }}.go
    echo 'package {{ name }}Model' > src/models/{{ name }}Model/{{ name }}.go

# scaffold a new controller (name is lowercase e.g. index)
@new-controller name:
    mkdir -pv src/controllers/{{ name }}Controller
    touch src/controllers/{{ name }}Controller/{{ name }}.go
    touch src/controllers/{{ name }}Controller/routes.go
    echo 'package {{ name }}Controller' > src/controllers/{{ name }}Controller/{{ name }}.go
    echo 'package {{ name }}Controller' > src/controllers/{{ name }}Controller/routes.go

[private]
make-autocap:
    #!/bin/bash
    mkif {{ AUTOCAP_BIN_PATH }} helpers/autocap/main.go -x 'just remake-autocap'

[private]
remake-autocap:
    #!/bin/bash
    set -e
    sudo rm -f {{ AUTOCAP_BIN_PATH }}
    go build -o {{ AUTOCAP_BIN_PATH }} helpers/autocap/main.go
    sudo chown root {{ AUTOCAP_BIN_PATH }}
    sudo chmod u+s {{ AUTOCAP_BIN_PATH }}

[private]
@redev: make
    {{ AUTOCAP_BIN_PATH }} {{ BIN_PATH }}
    {{ BIN_PATH }}

[private]
remake:
    #!/bin/bash
    set -e
    cd src
    go fmt ./...
    # We don't need the loopclosure check because of GOEXPERIMENT=loopvar.
    go vet -loopclosure=false ./...
    go build -o ../{{ BIN_PATH }} -trimpath -ldflags '-s -w'
