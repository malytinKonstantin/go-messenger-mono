FROM --platform=linux/amd64 golang:1.21 AS builder

WORKDIR /app

COPY . .
COPY ../proto ./proto

RUN apt-get update && apt-get install -y \
    build-essential \
    librdkafka-dev

RUN go mod download

RUN go build -o main ./api-gateway/cmd/server/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka1

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 3000
CMD ["./main"]