FROM golang

MAINTAINER allezsans <allez.sans.question@gmail.com>

WORKDIR ${GOPATH}/src/github.com/allezsans/yamato/go

COPY go/Gopkg.toml Gopkg.toml
COPY go/Gopkg.lock Gopkg.lock

RUN go get -u github.com/golang/dep/cmd/dep \
    && dep ensure --vendor-only
