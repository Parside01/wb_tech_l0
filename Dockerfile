FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0

RUN go build -o app ./cmd/main.go

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

RUN apk add --no-cache bash

CMD ["./app"]