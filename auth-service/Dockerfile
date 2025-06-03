# Stage 1: builder
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY shared ./shared
COPY auth-service/go.mod auth-service/go.sum ./auth-service/

WORKDIR /app/auth-service

RUN go mod download

COPY auth-service/. .

RUN go build -o main ./cmd/server

# Stage 2: final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service/main .
CMD ["./main"]