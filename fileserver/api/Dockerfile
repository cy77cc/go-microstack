FROM golang:1.25.5-alpine3.22 AS builder
ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app

COPY . .

WORKDIR /app/fileserver

RUN go mod tidy
RUN go build -ldflags="-s -w" -o /app/fileserver-api ./api/

FROM alpine:3.22
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/fileserver-api /app/fileserver-api
COPY fileserver/api/etc /app/etc
EXPOSE 41001
CMD ["./fileserver-api", "-f", "etc/fileserver.yaml"]
