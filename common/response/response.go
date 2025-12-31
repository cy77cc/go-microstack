package response

import (
	"net/http"

	"github.com/cy77cc/go-microstack/common/utils"
	"github.com/cy77cc/go-microstack/common/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Resp 定义统一的响应结构体
type Resp struct {
	Code      xcode.Xcode `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Error 返回错误响应
func Error(code xcode.Xcode, msg string) *Resp {
	return &Resp{
		Code:      code,
		Msg:       msg,
		Data:      nil,
		Timestamp: utils.GetTimestamp(),
	}
}

// Success 返回成功响应
func Success(data interface{}) *Resp {
	return &Resp{
		Code:      xcode.Success,
		Msg:       "success",
		Data:      data,
		Timestamp: utils.GetTimestamp(),
	}
}

// Response 统一处理 HTTP 响应
// w: http.ResponseWriter
// r: *http.Request
// resp: 成功时的数据
// err: 错误信息 (可能是 grpc error, custom error, or standard error)
func Response(w http.ResponseWriter, r *http.Request, resp interface{}, err error) {
	if err != nil {
		// 这里可以根据 err 类型进行更细致的处理
		// 简单处理：默认返回服务器错误，如果 err 是 xcode 类型则返回对应错误码
		// 实际项目中可能需要解析 grpc status code 等
		httpx.OkJson(w, Error(xcode.ServerError, err.Error()))
	} else {
		httpx.OkJson(w, Success(resp))
	}
}
