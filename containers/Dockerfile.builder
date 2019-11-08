FROM golang:1.11-alpine3.8

RUN apk add --no-cache ca-certificates cmake make g++ openssl-dev git

RUN go get golang.org/x/tools/cmd/cover github.com/mattn/goveralls github.com/DATA-DOG/godog/cmd/godog

WORKDIR /home/sehsyha/Code/crounch-back
COPY . .

ENV GO111MODULE=on
RUN go get