ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE AS builder 

RUN go build -o /go/bin/crounch-back -tags seabolt_static main.go

FROM alpine

COPY --from=builder /go/bin/crounch-back /go/bin/crounch-back 

EXPOSE 3000
ENTRYPOINT [ "/go/bin/crounch-back" ]