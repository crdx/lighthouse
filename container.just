IMAGE_NAME := 'lighthouse_app'
REMOTE_DIR := 'lighthouse'

# build the container
@build:
    docker-compose build

# start the container
@up: build
    docker-compose up

# stop the container
@down:
    docker-compose down

# start a shell in the app container
@shell: build
    docker-compose run app bash || true
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
