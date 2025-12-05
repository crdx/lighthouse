FROM golang:1.25.1-alpine3.22 AS build
# https://hub.docker.com/_/golang

RUN apk add --no-cache \
    build-base \
    git \
    libpcap-dev \
    libcap-utils

WORKDIR /build
COPY go.sum go.mod .
RUN go mod download

# Build.
COPY . .
RUN go build -o lighthouse -trimpath -ldflags '-s -w' ./cmd/lighthouse && \
    setcap cap_net_raw+eip lighthouse

# ——————————————————————————————————————————————————————————————————————————————————————————————————
FROM alpine:3.22.1
# https://hub.docker.com/_/alpine

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Europe/London /etc/localtime && \
    echo 'Europe/London' > /etc/timezone && \
    apk del tzdata

RUN apk add --no-cache \
    bash \
    libpcap-dev \
    curl \
    tzdata

RUN addgroup -g 1000 anon && \
    adduser -G anon -D -u 1000 anon

WORKDIR /app
COPY --from=build /build/lighthouse /init

# This needs to be a script within the container because we need access to $PORT.
RUN echo 'curl -sSf http://localhost:$PORT/health' >> healthcheck && \
    chmod +x healthcheck

USER anon

CMD ["/init"]
