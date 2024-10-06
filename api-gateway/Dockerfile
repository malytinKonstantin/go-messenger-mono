# FROM golang:1.23.1-alpine AS builder
# WORKDIR /app
# COPY . .
# RUN cd api-gateway && go mod download
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=builder /app/main .
# COPY --from=builder /app/api-gateway/.env .
# EXPOSE 3000
# CMD ["./main"]

FROM --platform=linux/amd64 golang:1.21 AS builder

WORKDIR /app
COPY . .
COPY ../proto /app/proto
COPY api-gateway/.env /app/.env

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
COPY --from=builder /app/.env ./.env

EXPOSE 3000
CMD ["./main"]