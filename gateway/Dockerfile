FROM golang:1.25.5-alpine3.22 AS builder

WORKDIR /app
COPY . .

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app/gateway

RUN go mod tidy

RUN go build -o /app/main cmd/gateway/main.go

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /app/main .

COPY gateway/etc ./etc

EXPOSE 41000

CMD ["./main", "-config", "etc/config.yaml"]