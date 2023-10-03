FROM golang:1.21.1-alpine3.18 AS build
# https://hub.docker.com/_/golang

RUN apk add --no-cache \
    build-base \
    git \
    libpcap-dev \
    libcap-utils \
    upx

WORKDIR /build/src
COPY src/go.sum src/go.mod ./

ENV GOEXPERIMENT=loopvar

# Fetch dependencies in a separate layer to speed up rebuilds.
ARG FORMAT='{{ if not .Main }}{{ .Path }}/...@{{ .Version }}{{ end }}'
RUN for PACKAGE in $(go list -m -f "$FORMAT" all); do go get $PACKAGE; done

# Build.
COPY src ./
RUN go build -o /build/dist/lighthouse -trimpath -ldflags '-s -w' && \
    upx /build/dist/lighthouse && \
    setcap cap_net_raw+eip /build/dist/lighthouse

# ——————————————————————————————————————————————————————————————————————————————————————————————————
FROM alpine:3.18.3
# https://hub.docker.com/_/alpine

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Europe/London /etc/localtime && \
    echo 'Europe/London' > /etc/timezone && \
    apk del tzdata

RUN apk add --no-cache \
    dumb-init \
    bash \
    libpcap-dev \
    curl \
    tzdata

RUN addgroup -g 1000 anon && \
    adduser -G anon -D -u 1000 anon

WORKDIR /app
COPY --from=build /build/dist/lighthouse lighthouse

# This needs to be a script within the container because we need access to $PORT.
RUN echo 'curl -sSf http://localhost:$PORT/health' >> healthcheck && \
    chmod +x healthcheck

USER anon

CMD ["dumb-init", "./lighthouse"]
