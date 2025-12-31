# User Center (用户中心)

## 简介
User Center 是一个基于 Go Zero 框架开发的用户中心微服务。它提供了用户管理、角色管理、权限管理（RBAC）以及认证鉴权功能。

## 功能特性

*   **用户管理**：注册、登录、用户信息更新、用户列表查询。
*   **角色管理**：创建角色、分配角色、角色列表。
*   **权限管理**：基于 RBAC 的权限控制，支持菜单、按钮、API 级别的权限粒度。
*   **认证鉴权**：基于 JWT 的认证机制，支持 Token 刷新和登出。
*   **微服务架构**：分离 API 层（HTTP）和 RPC 层（业务逻辑），易于扩展。

## 模块结构

*   `api/`: HTTP 接口层，处理外部请求。
*   `rpc/`: RPC 服务层，处理核心业务逻辑，直接与数据库交互。
*   `model/`: 数据库模型定义（Go Zero 生成）。

## API 接口概览

### 认证 (Auth)
*   `POST /usercenter/v1/auth/login`: 用户登录
*   `POST /usercenter/v1/auth/register`: 用户注册
*   `POST /usercenter/v1/auth/refresh`: 刷新 Token
*   `POST /usercenter/v1/auth/logout`: 用户登出

### 用户 (Users)
*   `POST /usercenter/v1/users`: 创建用户
*   `GET /usercenter/v1/users/:id`: 获取用户信息
*   `PUT /usercenter/v1/users/:id`: 更新用户信息
*   `DELETE /usercenter/v1/users/:id`: 删除用户
*   `GET /usercenter/v1/users`: 获取用户列表
*   `POST /usercenter/v1/users/:id/roles`: 分配角色
*   `GET /usercenter/v1/users/:id/roles`: 获取用户角色
*   `DELETE /usercenter/v1/users/:id/roles`: 移除用户角色

### 角色 (Roles)
*   `POST /usercenter/v1/roles`: 创建角色
*   `GET /usercenter/v1/roles`: 获取角色列表
*   `POST /usercenter/v1/roles/:id/permissions`: 分配权限

### 权限 (Permissions)
*   `POST /usercenter/v1/permissions`: 创建权限
*   `GET /usercenter/v1/permissions`: 获取权限列表

## 快速开始

### 前置要求
*   Go 1.18+
*   MySQL
*   Redis (可选，视配置而定)
*   Etcd (服务发现)

### 运行步骤

1.  **配置数据库**
    导入 `model/usercenter.sql` 和 `model/init_data.sql` 到 MySQL 数据库。

2.  **修改配置文件**
    *   修改 `rpc/etc/usercenter.yaml` 中的数据库连接串。
    *   修改 `api/etc/usercenter.yaml` 中的 RPC 服务配置。

3.  **启动 RPC 服务**
    ```bash
    cd rpc
    go run usercenter.go
    ```

4.  **启动 API 服务**
    ```bash
    cd api
    go run usercenter.go
    ```

5.  **访问接口**
    默认端口：API (8888), RPC (8080)
    Swagger 文档可参考 `.api` 文件描述。
