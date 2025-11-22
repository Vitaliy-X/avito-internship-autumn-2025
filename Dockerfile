# STAGE 1: build
# чтобы показать, что я знаю multi-stage build :)
FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pr-service ./cmd/pr-service

# STAGE 2: production
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/pr-service /app/pr-service

EXPOSE 8080

CMD ["./pr-service"]
