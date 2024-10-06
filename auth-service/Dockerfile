FROM --platform=linux/amd64 golang:1.21 AS builder

WORKDIR /app
COPY . .
COPY ../proto /app/proto
COPY auth-service/db /app/db
COPY auth-service/init.sh /app/init.sh
COPY auth-service/.env /app/.env

RUN apt-get update && apt-get install -y \
    build-essential \
    librdkafka-dev

RUN go mod download

RUN go build -o main ./auth-service/cmd/server/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka1

WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/db ./db
COPY --from=builder /app/init.sh ./init.sh
COPY --from=builder /app/.env ./.env

EXPOSE 50051 3001
RUN chmod +x ./init.sh
ENTRYPOINT ["./init.sh"]
CMD ["./main"]