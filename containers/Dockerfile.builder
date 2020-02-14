FROM golang:1.13-alpine

RUN apk add --no-cache ca-certificates cmake make g++ openssl-dev git

RUN go get golang.org/x/tools/cmd/cover github.com/mattn/goveralls github.com/cucumber/godog/cmd/godog

WORKDIR /home/sehsyha/Code/crounch-back
COPY . .

ENV GO111MODULE=on
