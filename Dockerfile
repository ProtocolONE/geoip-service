FROM golang:1.11.2-alpine AS builder
MAINTAINER Nikolay Bondarenko <nikolay.bondarenko@protocol.one>

RUN apk add bash ca-certificates git

ENV GO111MODULE=on
ENV MICRO_REGISTRY_ADDRESS=p1pay-consul
ENV MAXMIND_GEOIP_DB_PATH="/application/assets/maxmind/GeoLite2-City.mmdb"

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o $GOPATH/bin/geoip-service .

ENTRYPOINT $GOPATH/bin/geoip-service
