FROM golang:1.12.7-alpine@sha256:87e527712342efdb8ec5ddf2d57e87de7bd4d2fedf9f6f3547ee5768bb3c43ff AS builder

WORKDIR /build

RUN apk --no-cache add ca-certificates && update-ca-certificates

ADD . .

RUN GOFLAGS=-mod=vendor GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
    -a -ldflags '-extldflags "-static"' -o gitlab-project-settings .

FROM alpine:3.10.1@sha256:6a92cd1fcdc8d8cdec60f33dda4db2cb1fcdcacf3410a8e05b3741f44a9b5998

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /build/gitlab-project-settings   /bin/gitlab-project-settings
