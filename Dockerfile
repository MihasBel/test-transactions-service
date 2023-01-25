# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.1-alpine3.16 AS builder
WORKDIR /build

COPY . ./

RUN go mod download

RUN go build \
       			-o bin/app/test-transaction-service \
       			cmd/http/main.go

## Deploy
FROM alpine:3.16

WORKDIR /


COPY --from=builder /build/bin/app/test-transaction-service ./test-transaction-service
COPY configs/local/env.json ./

EXPOSE 8090



ENTRYPOINT ["/test-transactions"]
