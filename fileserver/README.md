# File Server (文件服务)

## 简介
File Server 是一个基于 Go Zero 框架开发的文件存储微服务。它提供了统一的文件上传、下载、管理接口，底层支持多种存储后端（目前主要支持 MinIO）。支持单文件上传和分片上传（Multipart Upload），适用于大文件传输场景。

## 功能特性

*   **单文件上传**：适用于小文件，直接上传。
*   **分片上传**：适用于大文件，支持断点续传、并发上传。
    *   初始化分片上传
    *   上传分片
    *   完成分片上传
    *   终止分片上传
*   **文件元数据管理**：记录文件大小、Hash、类型、上传者等信息。
*   **秒传**：基于文件 Hash 的重复文件检测，避免重复存储。
*   **多存储后端**：设计上支持多种后端（如 S3, MinIO, Local），目前实现对接 MinIO。

## 模块结构

*   `api/`: HTTP 接口层，处理外部请求，流式传输文件。
*   `rpc/`: RPC 服务层，处理文件元数据存储、与 MinIO 交互。
*   `model/`: 数据库模型定义。

## API 接口概览

### 文件操作 (Files)
*   `POST /fileserver/v1/files`: 单文件上传
*   `GET /fileserver/v1/files/:fileId/meta`: 获取文件元数据
*   `GET /fileserver/v1/files/:fileId/url`: 获取文件下载链接（预签名 URL）

### 分片上传 (Uploads)
*   `POST /fileserver/v1/uploads`: 初始化分片上传 (Initiate)
*   `PUT /fileserver/v1/uploads/:uploadId/parts/:partNumber`: 上传分片 (Upload Part)
*   `POST /fileserver/v1/uploads/:uploadId/complete`: 完成分片上传 (Complete)
*   `DELETE /fileserver/v1/uploads/:uploadId`: 终止分片上传 (Abort)

### 存储桶 (Buckets)
*   `POST /fileserver/v1/buckets`: 创建存储桶

## 分片上传流程

1.  **前端计算文件 Hash** (可选，用于秒传)。
2.  **调用 `POST /uploads`**：传入文件名、大小、Hash。
    *   后端调用 MinIO `InitiateMultipartUpload`，返回 `uploadId`。
3.  **前端切分文件**，并发调用 **`PUT /uploads/:uploadId/parts/:partNumber`**。
    *   每个分片上传到 MinIO。
4.  **前端调用 `POST /uploads/:uploadId/complete`**。
    *   后端调用 MinIO `CompleteMultipartUpload` 合并文件。
    *   后端保存文件记录到数据库。

## 快速开始

### 前置要求
*   Go 1.18+
*   MySQL
*   MinIO Server
*   Etcd

### 运行步骤

1.  **配置数据库**
    导入 `model/fileserver.sql` 到 MySQL。

2.  **启动 MinIO**
    确保 MinIO 运行并创建好 AccessKey/SecretKey。

3.  **修改配置文件**
    *   `rpc/etc/fileserver.yaml`: 配置 MySQL, MinIO 连接信息。
    *   `api/etc/fileserver.yaml`: 配置 RPC 服务发现。

4.  **启动 RPC 服务**
    ```bash
    cd rpc
    go run fileserver.go
    ```

5.  **启动 API 服务**
    ```bash
    cd api
    go run fileserver.go
    ```
