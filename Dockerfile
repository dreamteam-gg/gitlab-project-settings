FROM golang:1.12.7-alpine@sha256:87e527712342efdb8ec5ddf2d57e87de7bd4d2fedf9f6f3547ee5768bb3c43ff AS builder

WORKDIR /build

ADD . .

RUN GOFLAGS=-mod=vendor GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
    -a -ldflags '-extldflags "-static"' -o gitlab-project-settings .

FROM gcr.io/distroless/base@sha256:e37cf3289c1332c5123cbf419a1657c8dad0811f2f8572433b668e13747718f8

COPY --from=builder /build/gitlab-project-settings   /bin/gitlab-project-settings
