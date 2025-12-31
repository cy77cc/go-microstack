# Gateway Service

hioshop_ms 的 API 网关服务，基于 Gin 框架开发，提供统一的流量入口、路由转发、服务发现、负载均衡、熔断限流等功能。

## 功能特性

-   **服务发现**: 集成 Nacos 服务发现，自动感知后端服务实例。
-   **反向代理**: 基于 `httputil.ReverseProxy` 实现 HTTP 反向代理。
-   **负载均衡**: 支持轮询 (RoundRobin) 和随机 (Random) 策略。
-   **路由转发**: 支持基于路径前缀的动态路由配置。
-   **熔断降级**: 针对特定路由配置熔断策略 (Circuit Breaker)。
-   **流量控制**: 支持单机令牌桶限流 (Token Bucket) 和分布式限流 (Redis)。
-   **配置管理**: 支持本地 YAML 配置文件和 Nacos 动态配置加载。
-   **安全增强**: 自动注入 Request-ID、内部调用签名 (HMAC)、用户身份透传。

## 目录结构

```
gateway/
├── cmd/
│   └── gateway/
│       └── main.go           # 程序入口
├── config/                   # 配置定义与加载
├── etc/                      # 配置文件模板
├── internal/
│   ├── middleware/           # 中间件 (限流、熔断)
│   ├── proxy/                # 核心代理逻辑
│   └── router/               # 路由注册与管理
├── pkg/
│   └── loadbalance/          # 负载均衡算法
└── README.md
```

## 快速开始

### 1. 配置

在 `etc/config.yaml` 中配置基本服务信息和 Nacos 地址：

```yaml
server:
  name: gateway
  host: 0.0.0.0
  port: 8888

# Nacos 配置通过环境变量或代码硬编码加载，也可在此扩展
```

在 `etc/gateway-router.json` 中配置路由规则：

```json
[
  {
    "path_prefix": "/api/user",
    "service": "usercenter-api",
    "strip_prefix": "/api/user",
    "circuit": {
      "max_failures": 10,
      "open_seconds": 10
    },
    "rate_limit": {
      "qps": 100,
      "burst": 200
    }
  }
]
```

### 2. 启动

```bash
go run cmd/gateway/main.go
```

## 核心组件说明

### 代理处理器 (Proxy Handler)

位于 `internal/proxy/http_proxy.go`，负责：
1.  解析请求路径和服务名。
2.  通过服务发现获取可用实例列表。
3.  使用负载均衡器选择目标实例。
4.  重写 URL 路径（支持 `strip_prefix`）。
5.  注入内部安全头 (`X-Request-Id`, `X-User-Id`, `X-Signature`)。
6.  转发请求并处理响应。

### 中间件 (Middleware)

位于 `internal/middleware`，采用按路由配置动态加载的方式：
-   **CircuitBreaker**: 当后端服务错误率过高时自动熔断，快速失败，保护系统。
-   **RateLimit**: 限制特定路由的 QPS，防止流量洪峰打挂服务。

### 负载均衡 (LoadBalance)

位于 `pkg/loadbalance`，实现了简单的客户端负载均衡：
-   `RoundRobin`: 轮询算法，平滑分配流量。
-   `Random`: 随机算法。
