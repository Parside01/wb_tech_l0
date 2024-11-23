FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

RUN apk add --no-cache bash

CMD ["/wait-for-it.sh", "kafka:9092", "--timeout=30", "--strict", "--", "./app"]
CMD ["./app"]