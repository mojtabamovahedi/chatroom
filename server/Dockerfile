FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go



FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.json .

ENV CONFIG_PATH=/config.json

EXPOSE 8080
CMD ["./main"]
