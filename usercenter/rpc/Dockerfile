FROM golang:1.25.3-alpine3.22 AS builder
ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct

RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /app

COPY . .

WORKDIR /app/usercenter

RUN go mod tidy

RUN go build -ldflags="-s -w" -o /app/usercenter-rpc ./rpc/usercenter.go

FROM alpine:3.22
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/usercenter-rpc /app/usercenter-rpc
COPY usercenter/rpc/etc /app/etc
EXPOSE 40011
CMD ["./usercenter-rpc", "-f", "etc/usercenter.yaml"]