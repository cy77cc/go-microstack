# Common Library

hioshop_ms 项目的通用工具库，包含常用的加密、JWT、响应封装、错误码定义等模块。

## 模块说明

### 1. cryptx (加密模块)
提供基于 `scrypt` 的密码加密与验证功能。
- `PasswordEncrypt`: 密码加密
- `PasswordVerify`: 密码验证

### 2. jwtx (JWT 模块)
提供 JWT Token 的生成功能。
- `GetToken`: 生成标准 Token
- `GetTokenWithClaims`: 生成带自定义载荷的 Token

### 3. response (响应模块)
统一的 HTTP API 响应封装，规范化返回格式。
- `Success`: 成功响应
- `Error`: 错误响应
- `Response`: 自动处理响应（支持 error 类型判断）

### 4. xcode (错误码模块)
定义全局统一的业务错误码。
- `Success` (1000): 成功
- `ParamError` (2000): 参数错误
- `ServerError` (3000): 服务端错误
- ...

### 5. utils (工具模块)
常用辅助函数。
- `GenUUID`: 生成 UUID
- `GetFileExt`: 获取文件后缀
- `GetMachineIP`: 获取本机 IP
- `DeepCopy`: 结构体深拷贝

## 使用示例

### 响应封装
```go
import "github.com/cy77cc/go-microstack/common/response"

func MyHandler(w http.ResponseWriter, r *http.Request) {
    // ... 业务逻辑
    if err != nil {
        response.Response(w, r, nil, err)
        return
    }
    response.Response(w, r, data, nil)
}
```

### 错误码使用
```go
import "github.com/cy77cc/go-microstack/common/xcode"

return xcode.ParamError
```
