FROM golang:1.24.3-alpine3.21 AS build
# https://hub.docker.com/_/golang

RUN apk add --no-cache \
    build-base \
    git \
    libpcap-dev \
    libcap-utils \
    upx

WORKDIR /build
COPY go.sum go.mod .

# Fetch dependencies in a separate layer to speed up rebuilds.
ARG FORMAT='{{ if not .Main }}{{ .Path }}/...@{{ .Version }}{{ end }}'
RUN for PACKAGE in $(go list -m -f "$FORMAT" all); do go get $PACKAGE; done

# Build.
COPY . .
RUN go build -o lighthouse -trimpath -ldflags '-s -w' ./cmd/lighthouse && \
    upx lighthouse && \
    setcap cap_net_raw+eip lighthouse

# ——————————————————————————————————————————————————————————————————————————————————————————————————
FROM alpine:3.21.3
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
COPY --from=build /build/lighthouse lighthouse

# This needs to be a script within the container because we need access to $PORT.
RUN echo 'curl -sSf http://localhost:$PORT/health' >> healthcheck && \
    chmod +x healthcheck

USER anon

CMD ["./lighthouse"]
