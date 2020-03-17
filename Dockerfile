FROM golang:1.14.0-alpine3.11@sha256:6578dc0c1bde86ccef90e23da3cdaa77fe9208d23c1bb31d942c8b663a519fa5 AS builder

WORKDIR /build

RUN apk --no-cache add ca-certificates && update-ca-certificates

ADD . .

RUN GOFLAGS=-mod=vendor GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
    -a -ldflags '-extldflags "-static"' -o gitlab-project-settings .

FROM gcr.io/distroless/base:nonroot@sha256:54c459100e9d420e023b0aecc43f7010d2731b6163dd8e060906e2dec4c59890

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /build/gitlab-project-settings /bin/gitlab-project-settings
