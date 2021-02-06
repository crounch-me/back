ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE AS builder

RUN go build -o /go/bin/crounch-back main.go
# COPY init-db /go/bin/init-db

FROM alpine

COPY --from=builder /go/bin/crounch-back /go/bin/crounch-back
# COPY --from=builder /go/bin/init-db /init-db

EXPOSE 3000

ENTRYPOINT [ "/go/bin/crounch-back" , "serve", "--db-schema", "public", "--db-connection-uri", "postgresql://postgres:secretpassword@database/postgres?sslmode=disable"]
