FROM golang:1.22-alpine AS build_base

RUN mkdir -p /etc/ssl/certs/ && update-ca-certificates && apk add --no-cache git

WORKDIR /go/src/app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./main.go

FROM phusion/baseimage:focal-1.2.0

WORKDIR /app

RUN apt-get install -y ca-certificates
COPY --from=build_base /go/src/app/main /app/main

EXPOSE 8080

CMD ["./main"]