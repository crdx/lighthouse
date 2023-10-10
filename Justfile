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
    watchexec -w src -i src/conf/models.go -r 'just redev'

# start development with debug logging
@devd:
    LIGHTHOUSE_DEBUG=1 just dev

# start development with verbose logging
@devv:
    LIGHTHOUSE_VERBOSE=1 just dev

# open the development site in the default browser
@open:
    xdg-open http://$HOST:$PORT

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
@deploy: lint build
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
    cd src && go test -p 1 -cover ./... {{ args }} | grep -vF '[no test files]' || true

# connect to the test db
@test-db:
    mariadb {{ DB_NAME }}_test

# run linter
@lint:
    # We don't need the loopclosure check because of GOEXPERIMENT=loopvar.
    cd src && go vet -loopclosure=false ./... && golangci-lint run

# fix lint issues
@fix:
    cd src && golangci-lint run --fix

# format code
@fmt:
    cd src && go fmt ./...

# show code sloc
@sloc:
    tokei -tGo,HTML,CSS,JavaScript \
        -e src/assets/alpine.min.js \
        -e src/assets/bulma.min.css

# drop the db
@drop-db:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

# connect to the dev db
@db:
    mariadb {{ DB_NAME }}

# pull down the live db
@pull:
    importdb -f --host s --local {{ DB_NAME }} --remote {{ DB_NAME }}

# scaffold a new model (name is lowercase e.g. thing)
@new-model name:
    mkdir -pv src/models/{{ name }}M
    touch src/models/{{ name }}M/{{ name }}.go
    echo 'package {{ name }}M' > src/models/{{ name }}M/{{ name }}.go

# scaffold a new controller (name is lowercase e.g. index)
@new-controller name:
    mkdir -pv src/controllers/{{ name }}Controller
    touch src/controllers/{{ name }}Controller/routes.go
    touch src/controllers/{{ name }}Controller/routes_test.go
    echo 'package {{ name }}Controller' > src/controllers/{{ name }}Controller/routes.go
    echo 'package {{ name }}Controller_test' > src/controllers/{{ name }}Controller/routes_test.go

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
@generate:
    cd src && go generate ./...

[private]
@remake: generate
    cd src && go build -o ../{{ BIN_PATH }} -trimpath -ldflags '-s -w'
