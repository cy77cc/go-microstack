# Go Microstack

## 项目简介
Go Microstack 是一个基于 [Go Zero](https://go-zero.dev/) 微服务框架构建的后端系统。它旨在提供一套标准、高效、可扩展的微服务基础组件，包含用户中心、文件服务等核心模块。

## 模块列表

| 模块名称 | 目录 | 描述 |
| :--- | :--- | :--- |
| **User Center** | [usercenter](./usercenter) | 用户管理、认证鉴权 (JWT)、RBAC 权限控制 |
| **File Server** | [fileserver](./fileserver) | 文件上传下载、分片上传、对接 MinIO |

## 技术栈

*   **Golang**: 核心开发语言
*   **Go Zero**: 微服务框架 (RPC, API, Model)
*   **gRPC / Protobuf**: 服务间通信
*   **MySQL**: 关系型数据库
*   **Redis**: 缓存与会话管理
*   **Etcd**: 服务注册与发现
*   **MinIO**: 对象存储服务

## 快速开始

### 环境准备
1.  安装 Go 1.18+
2.  安装 `goctl` 工具: `go install github.com/zeromicro/go-zero/tools/goctl@latest`
3.  准备依赖服务: MySQL, Redis, Etcd, MinIO (推荐使用 Docker Compose 部署)

### 部署流程

1.  **初始化数据库**
    *   执行 `usercenter/model/*.sql`
    *   执行 `fileserver/model/*.sql`

2.  **启动 User Center**
    *   进入 `usercenter/rpc` 修改配置并启动 `go run usercenter.go`
    *   进入 `usercenter/api` 修改配置并启动 `go run usercenter.go`

3.  **启动 File Server**
    *   进入 `fileserver/rpc` 修改配置并启动 `go run fileserver.go`
    *   进入 `fileserver/api` 修改配置并启动 `go run fileserver.go`

## 开发规范

*   **API 定义**: 使用 `.api` 文件定义 HTTP 接口，通过 `goctl api go` 生成代码。
*   **RPC 定义**: 使用 `.proto` 文件定义 RPC 接口，通过 `goctl rpc protoc` 生成代码。
*   **错误处理**: 统一使用 `xerr` 包（需自定义）或标准错误码返回。
*   **文档**: 每个模块维护自己的 README，API 文档由 Swagger 自动生成（需配置）。

## 联系方式
*   Author: zhangdongping
