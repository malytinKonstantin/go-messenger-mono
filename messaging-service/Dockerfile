FROM golang:1.21-bookworm AS builder

WORKDIR /app

COPY . .
COPY ../proto ./proto

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    librdkafka-dev && \
    rm -rf /var/lib/apt/lists/*

RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic -o main ./messaging-service/cmd/server/main.go

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    librdkafka1 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 50054 3004

CMD ["./main"]