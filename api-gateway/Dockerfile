FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY . .
RUN cd api-gateway && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/api-gateway/.env .
EXPOSE 3000
CMD ["./main"]