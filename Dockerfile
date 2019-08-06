FROM golang:1.12.6-alpine3.9 AS builder
MAINTAINER Nikolay Bondarenko <nikolay.bondarenko@protocol.one>

RUN apk add bash ca-certificates git

WORKDIR /application

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/geoip-service .

FROM alpine:3.9

RUN apk update && apk add ca-certificates tzdata && rm -rf /var/cache/apk/*

USER 2010:2010

WORKDIR /application

COPY --from=builder /application /application

ENV MAXMIND_GEOIP_DB_PATH="/application/assets/GeoLite2-City.mmdb"

ENTRYPOINT /application/bin/geoip-service
