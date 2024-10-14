FROM golang:1.21-bookworm AS builder

WORKDIR /app

COPY . .
COPY ../proto ./proto

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    librdkafka-dev \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./api-gateway/cmd/server/main.go

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    librdkafka1 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]