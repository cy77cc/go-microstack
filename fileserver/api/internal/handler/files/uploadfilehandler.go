package files

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic/fileserver"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadFileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		var data []byte
		// Try to read from multipart form file "file"
		uploadFile, _, err := r.FormFile("file")
		if err == nil && uploadFile != nil {
			defer uploadFile.Close()
			data, err = io.ReadAll(uploadFile)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		} else {
			// Fallback to raw body
			data, err = io.ReadAll(r.Body)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "fileData", data)
		// uid 已在签名中间件注入，若中间件未启用则兜底读取头
		if _, ok := ctx.Value("uid").(uint64); !ok {
			if hv := r.Header.Get("X-User-Id"); hv != "" {
				if v, err := strconv.ParseUint(hv, 10, 64); err == nil {
					ctx = context.WithValue(ctx, "uid", v)
				}
			}
		}
		l := fileserver.NewUploadFileLogic(ctx, svcCtx)
		resp, err := l.UploadFile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
