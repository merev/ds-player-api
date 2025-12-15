# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o player-api ./cmd/player-api

# Runtime stage
FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=build /app/player-api /usr/local/bin/player-api

EXPOSE 8080
ENV APP_PORT=8080

CMD ["player-api"]
