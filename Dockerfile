FROM golang:1.14.0-alpine3.11@sha256:6578dc0c1bde86ccef90e23da3cdaa77fe9208d23c1bb31d942c8b663a519fa5 AS builder

WORKDIR /build

RUN apk --no-cache add ca-certificates && update-ca-certificates

ADD . .

RUN GOFLAGS=-mod=vendor GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
    -a -ldflags '-extldflags "-static"' -o gitlab-project-settings .

FROM alpine:3.11.3@sha256:ab00606a42621fb68f2ed6ad3c88be54397f981a7b70a79db3d1172b11c4367d

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /build/gitlab-project-settings /bin/gitlab-project-settings
