# syntax=docker/dockerfile:1
FROM golang:1.18 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o packsolver main.go

# final image
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/packsolver .
COPY static ./static

EXPOSE 8080
CMD ["./packsolver"]