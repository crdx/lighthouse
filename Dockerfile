FROM golang AS build
# https://hub.docker.com/_/golang

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
        git \
        libcap2-bin \
        libpcap-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# Build.
COPY . .
RUN --mount=type=cache,id=lighthouse,target=/root/.cache/go-build \
    go build -o lighthouse -trimpath -ldflags '-s -w' ./cmd/lighthouse && \
    setcap cap_net_raw+eip lighthouse

# ——————————————————————————————————————————————————————————————————————————————————————————————————
FROM debian:trixie-slim
# https://hub.docker.com/_/debian

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        curl \
        libpcap0.8 \
        tzdata && \
    rm -rf /var/lib/apt/lists/* && \
    cp /usr/share/zoneinfo/Europe/London /etc/localtime && \
    echo 'Europe/London' > /etc/timezone

RUN groupadd -g 1000 anon && \
    useradd -u 1000 -g anon -s /bin/sh anon

WORKDIR /app
COPY --from=build /build/lighthouse /init

# This needs to be a script within the container because we need access to $PORT.
RUN printf '%s\n' \
    '#!/bin/sh' \
    'exec curl -sSf "http://localhost:$PORT/health"' > healthcheck && \
    chmod +x healthcheck

USER anon

CMD ["/init"]
