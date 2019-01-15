FROM golang:1.11.2-alpine AS builder
MAINTAINER Nikolay Bondarenko <nikolay.bondarenko@protocol.one>

RUN apk add bash ca-certificates git

WORKDIR /application

ENV GO111MODULE=on
ENV MAXMIND_GEOIP_DB_PATH="/application/assets/GeoLite2-City.mmdb"

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o $GOPATH/bin/geoip-service .

ENTRYPOINT $GOPATH/bin/geoip-service
