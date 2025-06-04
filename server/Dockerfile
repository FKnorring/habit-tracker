# Build stage
FROM golang:1.21-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -o main ./server

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/data

EXPOSE 8080

ENV PORT=8080
ENV DB_PATH=/app/data/habits.db

# Run the binary
CMD ["./main"] 