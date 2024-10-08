FROM golang:1.21 AS builder

WORKDIR /app

COPY . .
COPY ../proto ./proto
COPY ../auth-service/db ./db
COPY ../auth-service/init.sh ./init.sh

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    librdkafka-dev \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN go mod download

RUN CGO_ENABLED=1 go build -tags dynamic -o main ./auth-service/cmd/server/main.go

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    librdkafka1 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/db ./db
COPY --from=builder /app/init.sh ./init.sh

EXPOSE 50051 3001
RUN chmod +x ./init.sh
ENTRYPOINT ["./init.sh"]
CMD ["./main"]