FROM golang:1.11-alpine3.8

RUN apk add --no-cache ca-certificates cmake make g++ openssl-dev git curl pkgconfig
RUN git clone -b v1.7.0 https://github.com/neo4j-drivers/seabolt.git /seabolt 
WORKDIR /seabolt/build
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

RUN go get golang.org/x/tools/cmd/cover
RUN go get github.com/mattn/goveralls
RUN go get github.com/DATA-DOG/godog/cmd/godog

WORKDIR $GOPATH/src/github.com/Sehsyha/crounch-back
COPY . .

ENV GO111MODULE=on
RUN go get