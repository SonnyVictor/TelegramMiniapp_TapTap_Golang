# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o main ./cmd/app/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080

# Chạy ứng dụng
CMD ["./main"]
