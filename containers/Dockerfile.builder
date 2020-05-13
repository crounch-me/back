FROM golang:1.13-alpine

RUN apk add --no-cache ca-certificates cmake make g++ openssl-dev git

WORKDIR /home/sehsyha/Code/crounch-back
COPY . .

ENV GO111MODULE=on
