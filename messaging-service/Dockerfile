FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .
EXPOSE 50054 3004
CMD ["./main"]