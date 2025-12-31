FROM golang:1.25.3-alpine3.22 AS builder
ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/usercenter-api ./api/usercenter.go

FROM alpine:3.22
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/usercenter-api /app/usercenter-api
COPY api/etc /app/etc
EXPOSE 41001
CMD ["./usercenter-api", "-f", "etc/usercenter.yaml"]