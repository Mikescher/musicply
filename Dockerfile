

FROM golang:1-alpine AS builder

RUN apk add --no-cache tzdata ca-certificates openssl make git tar coreutils bash curl

COPY . /buildsrc

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.1

RUN cd /buildsrc && make clean && make build



FROM alpine:latest

RUN apk add --no-cache tzdata

COPY --from=builder   /buildsrc/_build/musicply  /app/server

RUN mkdir /data

WORKDIR /app

EXPOSE 80

CMD ["/app/server"]
