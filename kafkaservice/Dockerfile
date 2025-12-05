FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o kafkaservice .

FROM golang:1.21-alpine

WORKDIR /app

RUN mkdir -p /app/logs

COPY --from=builder /app/kafkaservice .

CMD ["./kafkaservice"]
