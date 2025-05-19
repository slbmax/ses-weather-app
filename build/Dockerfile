FROM golang:1.24.0-alpine as buildbase

RUN apk add build-base

WORKDIR /build

ENV GOOS="linux"

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build  -o /usr/local/bin/weather-app /build

FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/weather-app /usr/local/bin/weather-app

RUN apk add --no-cache ca-certificates

ENTRYPOINT ["weather-app"]