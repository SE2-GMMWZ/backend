FROM golang:1.23.5-alpine3.20 AS builder

WORKDIR /app

COPY ./src ./src
COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/main ./src/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /build/main /app/main
# COPY ./migrations /app/migrations

ENTRYPOINT ["/app/main"]
