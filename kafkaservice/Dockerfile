FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o kafkaservice .

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /app/logs

COPY --from=builder /app/kafkaservice .

CMD ["./kafkaservice"]
