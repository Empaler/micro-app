FROM golang:1.24-bookworm AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

RUN apt-get update && apt-get install -y gcc musl-dev && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12 && \
    swag init -g cmd/api/main.go -o docs

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/api ./cmd/api

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates libpq-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=movies \
    SERVER_PORT=8080

CMD ["./api"]
